package main

import "github.com/pokemonpower92/collagegenerator/internal/logger"

func main() {
	logger.NewRequestLogger().Info("Hello from metadata worker!")
}
