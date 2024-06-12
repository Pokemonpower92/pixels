package main

import (
	"github.com/pokemonpower92/imagesetservice/config"
	"github.com/pokemonpower92/imagesetservice/internal/listener"
)

func main() {
	config.LoadEnvironmentVariables()
	isl := listener.NewImageSetConsumer()
	isl.Consume()
}
