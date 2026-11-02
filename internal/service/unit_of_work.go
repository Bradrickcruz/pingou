package service

import (
	"context"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// UnitOfWork define a interface para gerenciamento de transações.
// Permite que operações de repositório sejam executadas dentro de uma transação.
type UnitOfWork interface {
	// Begin inicia uma nova transação.
	// Retorna erro se já houver uma transação ativa.
	Begin(ctx context.Context) error

	// Commit confirma a transação ativa.
	// Retorna erro se não houver transação ativa.
	Commit() error

	// Rollback desfaz a transação ativa.
	// Pode ser chamado mesmo sem transação ativa (no-op).
	Rollback() error

	// CheckRepo retorna o repositório de checks para uso dentro da transação.
	CheckRepo() CheckRepositoryTx

	// MonitorRepo retorna o repositório de monitors para uso dentro da transação.
	MonitorRepo() MonitorRepositoryTx

	// IncidentRepo retorna o repositório de incidents para uso dentro da transação.
	IncidentRepo() IncidentRepositoryTx
}

// CheckRepositoryTx define métodos de repositório que aceitam transação.
// Usado para operações dentro de UnitOfWork.
type CheckRepositoryTx interface {
	CreateWithTx(ctx context.Context, c *domain.Check) error
	Latest(ctx context.Context, monitorID string, n int) ([]*domain.Check, error)
}

// MonitorRepositoryTx define métodos de repositório que aceitam transação.
type MonitorRepositoryTx interface {
	UpdateWithTx(ctx context.Context, m *domain.Monitor) error
	// UpdateCertWithTx grava os campos de certificado se a URL ainda for url.
	// Devolve false quando a URL mudou e nada foi gravado.
	UpdateCertWithTx(ctx context.Context, id, url string, expiresAt, alerted30For, alerted7For *time.Time) (bool, error)
	FindByID(ctx context.Context, id string) (*domain.Monitor, error)
}

// IncidentRepositoryTx define métodos de repositório que aceitam transação.
type IncidentRepositoryTx interface {
	CreateWithTx(ctx context.Context, i *domain.Incident) error
	FindOpenByMonitor(ctx context.Context, monitorID string) (*domain.Incident, error)
	Close(ctx context.Context, id string, endedAt string, durationSeconds int) error
	IncrementFailureCount(ctx context.Context, monitorID string) error
}
