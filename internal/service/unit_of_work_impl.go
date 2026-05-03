package service

import (
	"context"
	"database/sql"

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
	return u.checkRepo
}

func (u *sqliteUnitOfWork) MonitorRepo() MonitorRepositoryTx {
	return u.monitorRepo
}

func (u *sqliteUnitOfWork) IncidentRepo() IncidentRepositoryTx {
	return u.incidentRepo
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