package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(sqlitePath, migrationsPath string) error {
	dbDir := filepath.Dir(sqlitePath)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return fmt.Errorf("create db dir: %w", err)
	}

	absMigrationsPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("resolve migrations path: %w", err)
	}

	absSQLitePath, err := filepath.Abs(sqlitePath)
	if err != nil {
		return fmt.Errorf("resolve sqlite path: %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s", absMigrationsPath)
	databaseURL := fmt.Sprintf("sqlite3://%s", absSQLitePath)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
