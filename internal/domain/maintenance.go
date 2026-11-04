package domain

import "time"

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
