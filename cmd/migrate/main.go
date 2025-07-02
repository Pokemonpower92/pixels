package main

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/database"
)

func main() {
	time.Sleep(1000 * time.Millisecond)
	log.Printf("Migrating database...")
	connString := config.ConnString()
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	if err := database.RunMigration(config, "internal/sqlc/migrations/postgres"); err != nil {
		panic(err)
	} else {
		log.Printf("Migration succeeded.")
	}
}
