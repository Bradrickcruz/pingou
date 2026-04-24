package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite não suporta writes concorrentes

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	slog.Info("database ready")
	return db, nil
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrationFiles)
	goose.SetLogger(goose.NopLogger()) // silencia logs do goose

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Up(db, "migrations")
}
