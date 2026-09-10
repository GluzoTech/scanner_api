package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// migrationsDirectory resolves where the .sql files live.
//
// Defaults to "migrations" relative to the working directory, which is what a
// local `go run ./cmd/server` from the project root finds. A deployment that
// installs the binary elsewhere sets MIGRATIONS_DIR to an absolute path, so
// the server stops depending on the directory it happened to be started from.
func migrationsDirectory() string {
	if dir := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR")); dir != "" {
		return dir
	}

	return "migrations"
}

func RunMigrations(db *sql.DB) error {
	migrationsDir := migrationsDirectory()

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		// Name the path actually looked at: the usual failure is an installed
		// binary started from a directory that has no migrations/, and the
		// relative path alone does not reveal that.
		return fmt.Errorf(
			"failed to read migrations directory %s "+
				"(set MIGRATIONS_DIR to override): %w",
			resolvePath(migrationsDir),
			err,
		)
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

// resolvePath makes a path absolute for error messages, falling back to the
// original when that is not possible.
func resolvePath(path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return absolute
}
