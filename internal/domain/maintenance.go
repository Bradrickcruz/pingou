package domain

import (
	"fmt"
	"time"
)

// MaintenanceWindow é uma janela de manutenção agendada para um monitor.
// Enquanto ativa, o monitor não é conferido, não abre incidente e não avisa.
type MaintenanceWindow struct {
	ID          string
	MonitorID   string
	StartsAt    time.Time
	EndsAt      time.Time
	Reason      *string
	StartedAt   *time.Time // efeitos do início já aplicados (tick)
	FinalizedAt *time.Time // efeitos do fim já aplicados (tick ou cancelamento)
	CreatedAt   time.Time
}

// Active diz se o monitor está em manutenção em now. Vale do início até o fim
// ser aplicado (finalized_at), e não só até ends_at: entre ends_at e o tick que
// aplica o fim, um check com falha abriria um incidente que o fim deixaria
// órfão ao zerar o estado. O atraso é de no máximo um tick.
func (w *MaintenanceWindow) Active(now time.Time) bool {
	return w.FinalizedAt == nil && !now.Before(w.StartsAt)
}

// Started diz se a janela já começou pelo relógio (starts_at <= now). Janela
// iniciada tem o horário travado: só o motivo pode mudar (D12).
func (w *MaintenanceWindow) Started(now time.Time) bool {
	return !now.Before(w.StartsAt)
}

// MaintenanceWindowView é uma janela viva com o nome do monitor (página de manutenções).
type MaintenanceWindowView struct {
	MaintenanceWindow
	MonitorName string
}

// MaintenanceGroupKey agrupa as janelas que viram um aviso só: mesmo início,
// fim e motivo (o banco tem uma janela por monitor e nenhum id de lote). Motivo
// ausente é diferente de motivo preenchido, mesmo vazio.
func MaintenanceGroupKey(w *MaintenanceWindow) string {
	reason := "-"
	if w.Reason != nil {
		reason = "+" + *w.Reason
	}
	return fmt.Sprintf("%s|%s|%s", fmtKeyTime(w.StartsAt), fmtKeyTime(w.EndsAt), reason)
}

func fmtKeyTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }
