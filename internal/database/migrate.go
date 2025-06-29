package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

func RunMigration(config *pgx.ConnConfig) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return fmt.Errorf("could not find project root")
		}
		wd = parent
	}
	migrations := &migrate.FileMigrationSource{
		Dir: filepath.Join(wd, "internal", "sqlc", "migrations", "postgres"),
	}
	db := stdlib.OpenDB(*config)
	defer db.Close()
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	return nil
}
