package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// maxNoticeMonitors: quantos nomes entram no texto do e-mail e do ntfy. Uma
// janela pode ter 100 monitores, e o ntfy transforma mensagem acima de 4 KB em
// anexo. O webhook leva a lista inteira.
const maxNoticeMonitors = 20

// NewMaintenanceEvents monta as linhas do outbox de um aviso de manutenção:
// um event_id para o grupo e uma linha por canal. applied são as janelas do
// grupo e monitors os monitores delas. Sem canais ou sem monitores não há linha.
//
// monitor_id da linha é o primeiro monitor do grupo (a coluna é NOT NULL; o
// worker não depende dela para enviar). O prazo conta a partir de now.
func NewMaintenanceEvents(eventType string, applied []*domain.MaintenanceWindow, monitors []*domain.Monitor, channels []string, now time.Time) ([]domain.OutboxEvent, error) {
	snap := domain.NewMaintenanceSnapshot(applied, monitors, now)
	if len(channels) == 0 || len(snap.Maintenance.Monitors) == 0 {
		return nil, nil
	}
	eventID, err := newUUIDv7()
	if err != nil {
		return nil, err
	}
	rows := make([]domain.OutboxEvent, 0, len(channels))
	for _, ch := range channels {
		rows = append(rows, domain.OutboxEvent{
			EventID:       eventID,
			EventType:     eventType,
			Channel:       ch,
			MonitorID:     snap.Maintenance.Monitors[0].ID,
			Snapshot:      snap,
			Status:        domain.OutboxPending,
			CreatedAt:     now,
			ExpiresAt:     snap.OccurredAt.Add(incidentNoticeTTL),
			NextAttemptAt: now,
		})
	}
	return rows, nil
}

// maintenanceLink é o link da página de manutenções, ou nil sem PINGOU_PUBLIC_URL.
func maintenanceLink(publicURL string) *string {
	if publicURL == "" {
		return nil
	}
	u := publicURL + "/maintenance"
	return &u
}

// maintenanceSnapshot devolve o grupo do evento; erro se o snapshot não é de manutenção.
func maintenanceSnapshot(ev domain.OutboxEvent) (*domain.SnapshotMaintenance, error) {
	if ev.Snapshot.Maintenance == nil {
		return nil, fmt.Errorf("event %q without maintenance data", ev.EventType)
	}
	return ev.Snapshot.Maintenance, nil
}

// subjectMaintenance: um monitor cita o nome; vários citam só a contagem.
func subjectMaintenance(eventType string, mt *domain.SnapshotMaintenance) string {
	verb := "começou"
	if eventType == domain.EventMaintenanceEnded {
		verb = "terminou"
	}
	if len(mt.Monitors) == 1 {
		return fmt.Sprintf("[Pingou] 🔧 Manutenção %s: %s", verb, oneLine(mt.Monitors[0].Name))
	}
	return fmt.Sprintf("[Pingou] 🔧 Manutenção %s em %d monitores", verb, len(mt.Monitors))
}

// maintenanceBody monta o corpo (e-mail e ntfy), com as datas no fuso loc. O
// aviso de fim só diz que a manutenção terminou: o estado do monitor volta pelo
// fluxo normal de queda e recuperação.
func maintenanceBody(eventType string, mt *domain.SnapshotMaintenance, publicURL string, loc *time.Location) string {
	var b strings.Builder
	line := func(label, value string) { fmt.Fprintf(&b, "%s: %s\n", label, value) }

	if len(mt.Monitors) == 1 {
		line("Monitor", oneLine(mt.Monitors[0].Name))
	} else {
		b.WriteString("Monitores:\n")
		for i, m := range mt.Monitors {
			if i == maxNoticeMonitors {
				fmt.Fprintf(&b, "- e mais %d\n", len(mt.Monitors)-maxNoticeMonitors)
				break
			}
			fmt.Fprintf(&b, "- %s\n", oneLine(m.Name))
		}
	}
	line("Início", formatMessageTime(mt.StartsAt, loc))
	line("Fim", formatMessageTime(mt.EndsAt, loc))
	if mt.Reason != nil && oneLine(*mt.Reason) != "" {
		line("Motivo", oneLine(*mt.Reason))
	}
	if u := maintenanceLink(publicURL); u != nil {
		line("Detalhes", *u)
	}
	if eventType == domain.EventMaintenanceEnded {
		b.WriteString("\nA manutenção terminou.\n")
	}
	return b.String()
}
