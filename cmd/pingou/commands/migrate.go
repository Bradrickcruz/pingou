package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// migrateCmd e o comando pai para migrate
var migrateCmd = &cobra.Command{
	Use:                 "migrate",
	Short:               "Gerencia migrations do banco de dados",
	Long:                "Executa migrations (up, down, status) no banco de dados",
	PersistentPreRunE:     requireKey,
	TraverseChildren:      true,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Executa migrations pendentes",
	Long:  "Aplica todas as migrations pendentes",
	RunE:  runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Reverte a ultima migration",
	Long:  "Reverte uma migration",
	RunE:  runMigrateDown,
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra o status das migrations",
	Long:  "Lista todas as migrations e seu estado atual",
	RunE:  runMigrateStatus,
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("WARNING: failed to load .env: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	db, err := database.OpenDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("WARNING: failed to load .env: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	db, err := database.OpenDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	defer db.Close()

	if err := database.Down(db); err != nil {
		return fmt.Errorf("migration down failed: %w", err)
	}

	fmt.Println("Migration reverted successfully")
	return nil
}

func runMigrateStatus(cmd *cobra.Command, args []string) error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("WARNING: failed to load .env: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	db, err := database.OpenDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	defer db.Close()

	migrations, err := database.ListMigrations()
	if err != nil {
		return fmt.Errorf("list migrations failed: %w", err)
	}

	fmt.Println("Available migrations:")
	for _, m := range migrations {
		fmt.Println(" -", m)
	}

	rows, err := db.Query("SELECT version_id, is_applied, tstamp FROM goose_db_version ORDER BY version_id")
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	fmt.Println("\nApplied migrations:")
	for rows.Next() {
		var versionID int
		var isApplied int
		var tstamp string
		if err := rows.Scan(&versionID, &isApplied, &tstamp); err != nil {
			continue
		}
		status := "DOWN"
		if isApplied == 1 {
			status = "UP"
		}
		fmt.Printf("  %d %s\n", versionID, status)
	}

	return nil
}