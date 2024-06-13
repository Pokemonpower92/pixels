package main

import (
	"github.com/pokemonpower92/imagesetservice/config"
	"github.com/pokemonpower92/imagesetservice/internal/consumer"
)

func main() {
	config.LoadEnvironmentVariables()
	isc := consumer.NewImageSetConsumer()
	isc.Consume()
}
