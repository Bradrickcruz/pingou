package service

import (
	"context"
	"database/sql"

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
	CreateWithTx(ctx context.Context, tx *sql.Tx, c *domain.Check) error
	FindByMonitor(ctx context.Context, monitorID string, limit, offset int) ([]*domain.Check, int, error)
}

// MonitorRepositoryTx define métodos de repositório que aceitam transação.
type MonitorRepositoryTx interface {
	UpdateWithTx(ctx context.Context, tx *sql.Tx, m *domain.Monitor) error
	FindByID(ctx context.Context, id string) (*domain.Monitor, error)
}

// IncidentRepositoryTx define métodos de repositório que aceitam transação.
type IncidentRepositoryTx interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, i *domain.Incident) error
	FindOpenByMonitor(ctx context.Context, monitorID string) (*domain.Incident, error)
	Close(ctx context.Context, id string, endedAt string, durationSeconds int) error
}