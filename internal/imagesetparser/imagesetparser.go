package imagesetparser

import (
	"log"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/consumer"
	"github.com/pokemonpower92/collagegenerator/internal/jobhandler"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
)

func init() {
	config.LoadEnvironmentVariables()
}

func instantiateJobHandler() *jobhandler.JobHandler {
	dbConfig := config.NewPostgresConfig()
	repo, err := repository.NewImageSetRepository(dbConfig)
	if err != nil {
		panic(err)
	}
	jobHandlerLogger := log.New(log.Writer(), "jobhandler: ", log.Flags())
	return jobhandler.NewJobHandler(repo, jobHandlerLogger)
}

func instantiateImageSetConsumer(jh *jobhandler.JobHandler) *consumer.ImageSetConsumer {
	consumerLogger := log.New(log.Writer(), "consumer: ", log.Flags())
	consumerConfig := config.NewConsumerConfig()
	return consumer.NewImageSetConsumer(jh, consumerLogger, consumerConfig)
}

func Start() {
	jh := instantiateJobHandler()
	isc := instantiateImageSetConsumer(jh)
	isc.Consume()
}
