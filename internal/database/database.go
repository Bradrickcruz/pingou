package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// OpenDB abre conexao com banco SEM executar migrations
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", sqliteDSN(dsn))
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite não suporta writes concorrentes

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	slog.Info("database opened")
	return db, nil
}

// Open abre banco e executa migrations (comportamento atual)
func Open(dsn string) (*sql.DB, error) {
	db, err := OpenDB(dsn)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	slog.Info("database ready")
	return db, nil
}

// RunMigrations executa migrations pendentes
func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrationFiles)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Up(db, "migrations")
}

// Down executa rollback de uma migration
func Down(db *sql.DB) error {
	goose.SetBaseFS(migrationFiles)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Down(db, "migrations")
}

// Status mostra status das migrations
func Status(db *sql.DB) error {
	goose.SetBaseFS(migrationFiles)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Status(db, "migrations")
}

// ListMigrations lista os arquivos de migration disponiveis
func ListMigrations() ([]string, error) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

func sqliteDSN(dsn string) string {
	base, rawQuery, ok := strings.Cut(dsn, "?")
	values, err := url.ParseQuery(rawQuery)
	if !ok || err != nil {
		values = url.Values{}
	}

	values.Set("_foreign_keys", "on")
	values.Set("_journal_mode", "WAL")
	values.Set("_busy_timeout", "5000")

	return base + "?" + values.Encode()
}
