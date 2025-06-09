package main

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/database"
	"github.com/pokemonpower92/pixels/internal/repository"
)

func main() {
	time.Sleep(1000 * time.Millisecond)
	log.Printf("Migrating database...")
	postgresConfig := config.NewPostgresConfig()
	connString := repository.GetConnectionString(postgresConfig)
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	if err := database.RunMigration(config); err != nil {
		panic(err)
	} else {
		log.Printf("Migration succeeded.")
	}
}
