package imagesetparser

import (
	"log"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/consumer"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/generator"
	"github.com/pokemonpower92/collagegenerator/internal/jobhandler"
	"github.com/pokemonpower92/collagegenerator/internal/queue"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
)

func Start() {
	config.LoadEnvironmentVariables()

	dbConfig := config.NewPostgresConfig()
	repo, err := repository.NewImageSetRepository(dbConfig)
	if err != nil {
		panic(err)
	}

	store := datastore.NewS3Store()

	generatorLogger := log.New(log.Writer(), "generator: ", log.Flags())
	g := generator.NewImageSetGenerator(generatorLogger, store)

	jobHandlerLogger := log.New(log.Writer(), "jobhandler: ", log.Flags())
	jh := jobhandler.NewISJobHandler(jobHandlerLogger, repo, g)

	queueConfig := config.NewRabbitMQConfig()
	queueLogger := log.New(log.Writer(), "queue: ", log.Flags())
	q := queue.NewImageSetQueue(queueLogger, queueConfig)

	consumerLogger := log.New(log.Writer(), "consumer: ", log.Flags())
	isc := consumer.NewImageSetConsumer(consumerLogger, jh, q)
	isc.Consume()
}
