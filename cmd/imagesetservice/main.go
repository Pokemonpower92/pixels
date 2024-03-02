package main

import (
	"log"
	"os"

	"github.com/pokemonpower92/imagesetservice/config"
	"github.com/pokemonpower92/imagesetservice/internal/listener"
)

func main() {
	config.LoadEnvironmentVariables()

	l := log.New(os.Stdout, "collageapi", log.LstdFlags)
	isl := listener.NewImageSetConsumer(l)

	isl.Consume()
}
