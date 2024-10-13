package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
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
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	log.Println("Connected to database")
	log.Println("Creating imagesets table")
	_, err = pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS imagesets (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			description TEXT NOT NULL,
            type TEXT NOT NULL,
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
            file_name TEXT NOT NULL,
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

	log.Println("Seeding...")
	imq := fmt.Sprintf(
		`INSERT INTO imagesets (
            name,
            description,
            type
    ) values (
            '%s',
            'A testing imageset',
            'stock'
    ) RETURNING id;`, uuid.New())
	var id int
	err = pool.QueryRow(context.Background(), imq).Scan(&id)
	if err != nil {
		panic(err)
	}
	log.Printf("Found id: %d\n", id)

	acq := fmt.Sprintf(
		`INSERT INTO average_colors (
            imageset_id,
            file_name,
            red,
            green,
            blue,
            alpha
    ) values (
            %d,
            'test',
            1,
            2,
            3,
            0
    );`, id)
	row, err := pool.Query(context.Background(), acq)
	if err != nil {
		panic(err)
	}
	defer row.Close()
}
