package main

import (
	"context"
	"database/sql"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/rubenv/sql-migrate"
)

func run() error {
	ctx := context.Background()

	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	// Set up the migration source
	migrations := &migrate.FileMigrationSource{
		Dir: "internal/sqlc/migrations/sqlite3", // Adjust this path as needed
	}

	// Run the migrations
	_, err = migrate.Exec(conn, "sqlite3", migrations, migrate.Up)

	if err != nil {
		conn.Close()
		return err
	}

	queries := sqlc.New(conn)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, sqlc.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio: sql.NullString{
			String: "Co-author of The C Programming Language and The Go Programming Language",
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
