package listener

import (
	"encoding/json"
	"log"

	"github.com/pokemonpower92/imagesetservice/config"
	"github.com/pokemonpower92/imagesetservice/internal/imageset"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ImageSetConsumer struct {
	l          *log.Logger
	config     *config.ConsumerConfig
	JobHandler *imageset.JobHandler
}

func NewImageSetConsumer(l *log.Logger) *ImageSetConsumer {
	return &ImageSetConsumer{
		l:          l,
		config:     config.NewConsumerConfig(),
		JobHandler: imageset.NewJobHandler(l),
	}
}

func (isc *ImageSetConsumer) Consume() {
	conn, err := amqp.Dial(isc.config.RabbitMQConfig.URI)
	if err != nil {
		isc.l.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		isc.l.Fatalf("Failed to open a channel: %s", err)
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
		isc.l.Fatalf("Failed to declare a queue: %s", err)
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
		isc.l.Fatalf("Failed to register a consumer: %s", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			isc.l.Printf("Received a job: %s", d.Body)

			var job map[string]interface{}
			err := json.Unmarshal(d.Body, &job)
			if err != nil {
				isc.l.Printf("Failed to decode job: %s", err)
			}

			decodedJob := imageset.NewJob(job)
			isc.JobHandler.HandleJob(decodedJob)
		}
	}()

	isc.l.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
