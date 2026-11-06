package service

import (
	"context"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// Atalhos para os testes dos canais: montam o evento do outbox a partir do
// monitor e do incidente e chamam Send.

func testEvent(eventType string, snap domain.EventSnapshot) domain.OutboxEvent {
	id, _ := newUUIDv7()
	return domain.OutboxEvent{EventID: id, EventType: eventType, Channel: "teste", MonitorID: snap.Monitor.ID, Snapshot: snap}
}

func notifyDown(ctx context.Context, c Channel, m *domain.Monitor, i *domain.Incident) error {
	return c.Send(ctx, testEvent(domain.EventDown, domain.NewDownSnapshot(m, i)))
}

func notifyRecovery(ctx context.Context, c Channel, m *domain.Monitor, i *domain.Incident) error {
	return c.Send(ctx, testEvent(domain.EventUp, domain.NewRecoverySnapshot(m, i, time.Now())))
}

func notifyCertExpiring(ctx context.Context, c Channel, m *domain.Monitor, a domain.CertAlert) error {
	return c.Send(ctx, testEvent(domain.EventCertExpiring, domain.NewCertSnapshot(m, a, time.Now())))
}
