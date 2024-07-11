package consumer

import (
	"encoding/json"
	"log"

	"github.com/pokemonpower92/collagegenerator/internal/job"
	"github.com/pokemonpower92/collagegenerator/internal/jobhandler"
	"github.com/pokemonpower92/collagegenerator/internal/queue"
)

type ImageSetConsumer struct {
	logger     *log.Logger
	jobHandler jobhandler.JobHandler
	queue      queue.Queue
}

func NewImageSetConsumer(
	jobhandler jobhandler.JobHandler,
	logger *log.Logger,
	queue queue.Queue) *ImageSetConsumer {
	return &ImageSetConsumer{
		logger:     logger,
		jobHandler: jobhandler,
		queue:      queue,
	}
}

func (isc *ImageSetConsumer) Consume() {
	isc.queue.Initialize()
	defer isc.queue.Close()

	go func() {
		for d := range isc.queue.Consume() {
			isc.logger.Printf("Received a job: %s", d.Body)

			var j map[string]interface{}
			err := json.Unmarshal(d.Body, &j)
			if err != nil {
				isc.logger.Printf("Failed to decode job: %s", err)
			}

			decodedJob := job.NewJob(j)
			err = isc.jobHandler.HandleJob(decodedJob)
			if err != nil {
				isc.logger.Printf("Job failed with error: %s", err)
			}
		}
	}()

	isc.logger.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}
