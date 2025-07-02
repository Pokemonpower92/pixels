package database

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

func RunMigration(
	config *pgx.ConnConfig,
	migrationPath string,
) error {
	migrations := &migrate.FileMigrationSource{
		Dir: migrationPath,
	}
	db := stdlib.OpenDB(*config)
	defer db.Close()
	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	return nil
}
