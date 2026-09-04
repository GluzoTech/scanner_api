package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func RunMigrations(db *sql.DB) error {
	migrationsDir := "migrations"

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		path := filepath.Join(migrationsDir, entry.Name())

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf(
				"failed to read migration %s: %w",
				entry.Name(),
				err,
			)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf(
				"failed to execute migration %s: %w",
				entry.Name(),
				err,
			)
		}
	}

	return nil
}
