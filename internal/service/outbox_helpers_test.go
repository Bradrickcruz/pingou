package service_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

// testChannels: um canal só, então cada aviso vira uma linha no outbox.
var testChannels = []string{domain.ChannelWebhook}

// newUoW monta a UnitOfWork de teste, com o outbox.
func newUoW(db *sql.DB) service.UnitOfWork {
	return service.NewUnitOfWork(db, repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db),
		repository.NewIncidentRepoTx(db), repository.NewOutboxRepoTx())
}

// pendingNotices devolve os avisos pending do outbox do tipo eventType ("" = todos),
// na ordem em que foram gravados.
func pendingNotices(t *testing.T, db *sql.DB, eventType string) []domain.OutboxEvent {
	t.Helper()
	all, err := repository.NewOutboxRepo(db).Due(t.Context(), time.Now().Add(100*365*24*time.Hour), 10000)
	if err != nil {
		t.Fatal(err)
	}
	var out []domain.OutboxEvent
	for _, ev := range all {
		if eventType == "" || ev.EventType == eventType {
			out = append(out, ev)
		}
	}
	return out
}
