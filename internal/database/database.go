package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

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

// RunMigrations executa migrations pendentes. Avisa no log antes e depois:
// num banco grande, criar um índice leva de segundos a minutos, e sem o aviso
// a subida parece travada.
func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrationFiles)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	current, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}
	pending, err := goose.CollectMigrations("migrations", current, goose.MaxVersion)
	if errors.Is(err, goose.ErrNoMigrationFiles) || (err == nil && len(pending) == 0) {
		return nil // banco em dia: o goose trata "nada depois da versão atual" como erro
	}
	if err != nil {
		return err
	}
	versions := make([]int64, len(pending))
	for i, m := range pending {
		versions[i] = m.Version
	}
	slog.Info("applying database migrations (may take a while on a large database)",
		"from", current, "pending", versions)

	start := time.Now()
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	slog.Info("database migrations applied",
		"version", versions[len(versions)-1],
		"took", time.Since(start).Round(time.Millisecond).String())
	return nil
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
