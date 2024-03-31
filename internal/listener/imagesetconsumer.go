package listener

import (
	"encoding/json"
	"fmt"
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

func NewImageSetConsumer() *ImageSetConsumer {
	return &ImageSetConsumer{
		l:          log.New(log.Writer(), "imagesetconsumer ", log.LstdFlags),
		config:     config.NewConsumerConfig(),
		JobHandler: imageset.NewJobHandler(),
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
