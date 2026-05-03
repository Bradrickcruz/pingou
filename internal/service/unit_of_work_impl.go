package service

import (
	"context"
	"database/sql"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
)

// sqliteUnitOfWork implementação de UnitOfWork para SQLite.
type sqliteUnitOfWork struct {
	db           *sql.DB
	checkRepo    *repository.CheckRepoTx
	monitorRepo  *repository.MonitorRepoTx
	incidentRepo *repository.IncidentRepoTx
	tx           *sql.Tx
}

// NewUnitOfWork cria uma nova instância de UnitOfWork.
// Recebe a conexão DB e instâncias dos repositórios.
func NewUnitOfWork(db *sql.DB, checkRepo *repository.CheckRepoTx, monitorRepo *repository.MonitorRepoTx, incidentRepo *repository.IncidentRepoTx) UnitOfWork {
	return &sqliteUnitOfWork{
		db:           db,
		checkRepo:    checkRepo,
		monitorRepo:  monitorRepo,
		incidentRepo: incidentRepo,
	}
}

func (u *sqliteUnitOfWork) Begin(ctx context.Context) error {
	if u.tx != nil {
		return ErrTransactionAlreadyActive
	}
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.tx = tx
	return nil
}

func (u *sqliteUnitOfWork) Commit() error {
	if u.tx == nil {
		return ErrNoActiveTransaction
	}
	err := u.tx.Commit()
	u.tx = nil
	return err
}

func (u *sqliteUnitOfWork) Rollback() error {
	if u.tx == nil {
		return nil // no-op se não há transação ativa
	}
	err := u.tx.Rollback()
	u.tx = nil
	return err
}

func (u *sqliteUnitOfWork) CheckRepo() CheckRepositoryTx {
	return &checkRepoWithTx{u.checkRepo, u.tx}
}

func (u *sqliteUnitOfWork) MonitorRepo() MonitorRepositoryTx {
	return &monitorRepoWithTx{u.monitorRepo, u.tx}
}

func (u *sqliteUnitOfWork) IncidentRepo() IncidentRepositoryTx {
	return &incidentRepoWithTx{u.incidentRepo, u.tx}
}

// checkRepoWithTx envolve CheckRepoTx com uma transação ativa
type checkRepoWithTx struct {
	repo *repository.CheckRepoTx
	tx   *sql.Tx
}

func (r *checkRepoWithTx) CreateWithTx(ctx context.Context, c *domain.Check) error {
	return r.repo.CreateWithTx(ctx, r.tx, c)
}

func (r *checkRepoWithTx) FindByMonitor(ctx context.Context, monitorID string, limit, offset int) ([]*domain.Check, int, error) {
	return r.repo.FindByMonitor(ctx, monitorID, limit, offset)
}

// monitorRepoWithTx envolve MonitorRepoTx com uma transação ativa
type monitorRepoWithTx struct {
	repo *repository.MonitorRepoTx
	tx   *sql.Tx
}

func (r *monitorRepoWithTx) UpdateWithTx(ctx context.Context, m *domain.Monitor) error {
	return r.repo.UpdateWithTx(ctx, r.tx, m)
}

func (r *monitorRepoWithTx) FindByID(ctx context.Context, id string) (*domain.Monitor, error) {
	return r.repo.FindByID(ctx, id)
}

// incidentRepoWithTx envolve IncidentRepoTx com uma transação ativa
type incidentRepoWithTx struct {
	repo *repository.IncidentRepoTx
	tx   *sql.Tx
}

func (r *incidentRepoWithTx) CreateWithTx(ctx context.Context, i *domain.Incident) error {
	return r.repo.CreateWithTx(ctx, r.tx, i)
}

func (r *incidentRepoWithTx) FindOpenByMonitor(ctx context.Context, monitorID string) (*domain.Incident, error) {
	return r.repo.FindOpenByMonitor(ctx, monitorID)
}

func (r *incidentRepoWithTx) Close(ctx context.Context, id string, endedAt string, durationSeconds int) error {
	return r.repo.Close(ctx, id, endedAt, durationSeconds)
}

// Erros específicos de UnitOfWork.
var (
	ErrTransactionAlreadyActive = &UoWError{Message: "transaction already active"}
	ErrNoActiveTransaction      = &UoWError{Message: "no active transaction"}
)

type UoWError struct {
	Message string
}

func (e *UoWError) Error() string {
	return e.Message
}
