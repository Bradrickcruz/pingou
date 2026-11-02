package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type StateMachine struct {
	uow      UnitOfWork
	notifier Notifier
	clock    func() time.Time
}

func NewStateMachine(uow UnitOfWork, notifier Notifier) *StateMachine {
	return &StateMachine{
		uow:      uow,
		notifier: notifier,
		clock:    time.Now,
	}
}

// WithClock troca o relógio (testes simulam dias de checks sem esperar).
func (sm *StateMachine) WithClock(clock func() time.Time) *StateMachine {
	sm.clock = clock
	return sm
}

// Process recebe o resultado de um check e aplica as transições de estado.
// Todas as operações de persistência são executadas dentro de uma transação.
func (sm *StateMachine) Process(ctx context.Context, m *domain.Monitor, result domain.CheckResult) error {
	now := sm.clock().UTC()

	// Inicia transação
	if err := sm.uow.Begin(ctx); err != nil {
		return err
	}

	// Garante rollback em caso de erro
	success := false
	defer func() {
		if !success {
			_ = sm.uow.Rollback()
		}
	}()

	// 1. persiste o check dentro da transação
	check := &domain.Check{
		MonitorID:    m.ID,
		Success:      result.Success,
		StatusCode:   result.StatusCode,
		LatencyMs:    result.LatencyMs,
		ErrorMessage: result.ErrorMessage,
		CheckedAt:    now,
	}
	if err := sm.uow.CheckRepo().CreateWithTx(ctx, check); err != nil {
		return err
	}

	// 2. aplica transição de estado
	prevState := m.CurrentState

	var notify func()
	var err error
	if result.Success {
		notify, err = sm.handleSuccess(ctx, m, prevState, now)
	} else {
		notify, err = sm.handleFailure(ctx, m, prevState, result, now)
	}
	if err != nil {
		return err
	}

	// 3. certificado SSL: data de expiração e aviso de 30/7 dias
	certNotify, err := sm.processCert(ctx, m, result, now)
	if err != nil {
		return err
	}

	// Commit da transação antes de notificar
	if err := sm.uow.Commit(); err != nil {
		return err
	}

	success = true

	// Notificações APÓS commit: o notifier lê settings no banco, e com
	// MaxOpenConns=1 isso travaria esperando a conexão presa na transação
	if notify != nil {
		notify()
	}
	if certNotify != nil {
		certNotify()
	}

	return nil
}

// processCert atualiza a data do certificado e decide o aviso de expiração,
// dentro da transação do check. Só escreve quando algo muda (data nova ou marca
// nova): na maioria dos checks não há escrita extra. Devolve o aviso a enviar
// após o commit (ou nil).
//
// Check sem certificado (falha de conexão, http://) não apaga a data guardada,
// e o aviso ainda é avaliado com ela: um site fora do ar não adia o aviso.
func (sm *StateMachine) processCert(ctx context.Context, m *domain.Monitor, result domain.CheckResult, now time.Time) (func(), error) {
	if m.Type != domain.TypeHTTP {
		return nil, nil
	}

	changed := false
	// CertURL diferente: o usuário trocou a URL enquanto o check rodava, e o
	// certificado lido é o do endereço antigo
	if result.CertNotAfter != nil && result.CertURL == m.URL {
		notAfter := result.CertNotAfter.UTC().Truncate(time.Second) // precisão do banco
		if m.CertExpiresAt == nil || !m.CertExpiresAt.Equal(notAfter) {
			m.CertExpiresAt = &notAfter
			changed = true
		}
	}

	var alert *domain.CertAlert
	if m.CertAlertEnabled && m.CertExpiresAt != nil {
		d := EvaluateCertAlert(now, *m.CertExpiresAt, m.CertAlerted30For, m.CertAlerted7For)
		expiresAt := *m.CertExpiresAt
		if d.Mark30 {
			m.CertAlerted30For = &expiresAt
		}
		if d.Mark7 {
			m.CertAlerted7For = &expiresAt
		}
		changed = changed || d.Mark30 || d.Mark7
		alert = d.Alert
	}

	if !changed {
		return nil, nil
	}
	// guardado pela URL: uma edição feita entre a leitura do monitor e esta
	// transação zera os campos, e o certificado antigo não volta
	saved, err := sm.uow.MonitorRepo().UpdateCertWithTx(ctx, m.ID, m.URL, m.CertExpiresAt, m.CertAlerted30For, m.CertAlerted7For)
	if err != nil {
		return nil, err
	}
	if !saved || alert == nil {
		return nil, nil
	}

	a := *alert
	slog.Info("cert expiring alert", "monitor_id", m.ID, "name", m.Name,
		"threshold_days", a.ThresholdDays, "days_left", a.DaysLeft, "cert_expires_at", a.ExpiresAt.Format(time.RFC3339))
	return func() { sm.notifier.NotifyCertExpiring(ctx, m, a) }, nil
}

// handleSuccess aplica a transição e devolve a notificação a enviar após o commit (ou nil).
func (sm *StateMachine) handleSuccess(ctx context.Context, m *domain.Monitor, prevState domain.MonitorState, now time.Time) (func(), error) {
	m.CurrentState = domain.StateUp
	m.LastCheckedAt = &now
	m.UpdatedAt = now

	if err := sm.uow.MonitorRepo().UpdateWithTx(ctx, m); err != nil {
		return nil, err
	}

	// DOWN→UP: fecha incidente (dentro da transação)
	if prevState == domain.StateDown {
		incident, err := sm.uow.IncidentRepo().FindOpenByMonitor(ctx, m.ID)
		if err != nil {
			return nil, err
		}
		if incident != nil {
			endedAt := now.Format(time.RFC3339)
			durationSeconds := int(now.Sub(incident.StartedAt).Seconds())
			if durationSeconds < 0 {
				durationSeconds = 0
			}
			if err := sm.uow.IncidentRepo().Close(ctx, incident.ID, endedAt, durationSeconds); err != nil {
				return nil, err
			}
			d := int64(durationSeconds)
			incident.DurationSeconds = &d
			incident.EndedAt = &now
			slog.Info("monitor recovered", "monitor_id", m.ID, "name", m.Name)
			return func() { sm.notifier.NotifyRecovery(ctx, m, incident) }, nil
		}
	}
	// UNKNOWN→UP: não notifica (D na PRD)

	return nil, nil
}

// handleFailure aplica a transição e devolve a notificação a enviar após o commit (ou nil).
func (sm *StateMachine) handleFailure(ctx context.Context, m *domain.Monitor, prevState domain.MonitorState, result domain.CheckResult, now time.Time) (func(), error) {
	m.LastCheckedAt = &now
	m.UpdatedAt = now

	// conta falhas consecutivas recentes
	consecutiveFails, err := sm.countConsecutiveFails(ctx, m)
	if err != nil {
		return nil, err
	}

	// threshold ainda não atingido — mantém estado atual (se já estava DOWN,
	// por ex. threshold aumentado durante a queda, a falha conta no incidente)
	if consecutiveFails < effectiveThreshold(m) {
		if err := sm.uow.MonitorRepo().UpdateWithTx(ctx, m); err != nil {
			return nil, err
		}
		if prevState == domain.StateDown {
			return nil, sm.uow.IncidentRepo().IncrementFailureCount(ctx, m.ID)
		}
		return nil, nil
	}

	// threshold atingido — vai pra DOWN
	m.CurrentState = domain.StateDown
	if err := sm.uow.MonitorRepo().UpdateWithTx(ctx, m); err != nil {
		return nil, err
	}

	// só abre incidente na transição →DOWN (não repete se já estava DOWN)
	if prevState != domain.StateDown {
		id, err := newUUIDv7()
		if err != nil {
			return nil, err
		}
		// as falhas que atingiram o threshold já contam como "falhas seguidas"
		incident := &domain.Incident{
			ID:           id,
			MonitorID:    m.ID,
			StartedAt:    now,
			LastError:    result.ErrorMessage,
			FailureCount: int64(consecutiveFails),
		}
		if err := sm.uow.IncidentRepo().CreateWithTx(ctx, incident); err != nil {
			return nil, err
		}
		slog.Warn("monitor down", "monitor_id", m.ID, "name", m.Name, "error", result.ErrorMessage)
		return func() { sm.notifier.NotifyDown(ctx, m, incident) }, nil
	}

	// já estava DOWN: mais uma falha no incidente aberto (last_error segue
	// sendo o motivo da queda)
	return nil, sm.uow.IncidentRepo().IncrementFailureCount(ctx, m.ID)
}

// countConsecutiveFails conta checks recentes até encontrar um sucesso
// Lê dentro da transação (SQLite tem 1 conexão; ler fora dela trava)
func (sm *StateMachine) countConsecutiveFails(ctx context.Context, m *domain.Monitor) (int, error) {
	// busca os últimos N checks (N = failure_threshold + 1 para ter margem)
	checks, err := sm.uow.CheckRepo().Latest(ctx, m.ID, effectiveThreshold(m)+1)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, c := range checks {
		if !c.Success {
			count++
		} else {
			break // encontrou sucesso, para de contar
		}
	}
	return count, nil
}

// effectiveThreshold: heartbeat vira DOWN na primeira falha (prazo perdido);
// a tolerância (grace_seconds) já cumpre o papel do threshold.
func effectiveThreshold(m *domain.Monitor) int {
	if m.Type == domain.TypeHeartbeat {
		return 1
	}
	return m.FailureThreshold
}
