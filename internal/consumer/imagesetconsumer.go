package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pokemonpower92/imagesetservice/config"
	"github.com/pokemonpower92/imagesetservice/internal/job"
	"github.com/pokemonpower92/imagesetservice/internal/jobhandler"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ImageSetConsumer struct {
	logger     *log.Logger
	jobHandler *jobhandler.JobHandler
	config     *config.ConsumerConfig
}

func NewImageSetConsumer(
	jobhandler *jobhandler.JobHandler,
	logger *log.Logger,
	config *config.ConsumerConfig) *ImageSetConsumer {
	return &ImageSetConsumer{
		logger:     logger,
		jobHandler: jobhandler,
		config:     config,
	}
}

func (isc *ImageSetConsumer) Consume() {
	connString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		isc.config.RabbitMQConfig.User,
		isc.config.RabbitMQConfig.Password,
		isc.config.RabbitMQConfig.Host,
		isc.config.RabbitMQConfig.Port,
	)
	conn, err := amqp.Dial(connString)
	if err != nil {
		isc.logger.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		isc.logger.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		isc.config.Queue.Name,
		isc.config.Queue.Durable,
		isc.config.Queue.AutoDelete,
		isc.config.Queue.Exclusive,
		isc.config.Queue.NoWait,
		isc.config.Queue.Args,
	)
	if err != nil {
		isc.logger.Fatalf("Failed to declare a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		isc.logger.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
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
