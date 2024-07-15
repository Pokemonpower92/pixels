package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/collagegenerator/config"
)

func main() {
	config.LoadEnvironmentVariables()

	log := log.New(log.Writer(), "create_db: ", log.LstdFlags)

	config := config.NewPostgresConfig()
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
	log.Printf("Connecting to %s\n", connString)

	pool, err := pgxpool.New(
		context.Background(),
		connString,
	)
	defer pool.Close()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to database")

	log.Println("Creating imagesets table")
	_, err = pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS imagesets (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			description TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	)
	if err != nil {
		panic(err)
	}

	log.Println("Creating average_colors table")
	_, err = pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS average_colors (
			id SERIAL PRIMARY KEY,
			imageset_id INT NOT NULL,
			red INT NOT NULL,
			green INT NOT NULL,
			blue INT NOT NULL,
			alpha INT NOT NULL,
			FOREIGN KEY (imageset_id) REFERENCES imagesets(id)
		);`,
	)
	if err != nil {
		panic(err)
	}

	log.Println("Database initialized successfully")
}
