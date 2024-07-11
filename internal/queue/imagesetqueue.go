package queue

import (
	"fmt"
	"log"

	"github.com/pokemonpower92/collagegenerator/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ImageSetQueue struct {
	logger  *log.Logger
	config  *config.RabbitMQConfig
	channel *amqp.Channel
	queue   *amqp.Queue
	conn    *amqp.Connection
}

func NewImageSetQueue(logger *log.Logger, config *config.RabbitMQConfig) *ImageSetQueue {
	return &ImageSetQueue{
		logger: logger,
		config: config,
	}
}

func (isq *ImageSetQueue) Initialize() {
	connString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		isq.config.User,
		isq.config.Password,
		isq.config.Host,
		isq.config.Port,
	)

	conn, err := amqp.Dial(connString)
	if err != nil {
		isq.logger.Fatalf("Failed to connect to RabbitMQ: %s", err)
		panic(err)
	}
	isq.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		isq.logger.Fatalf("Failed to open a channel: %s", err)
		panic(err)
	}
	isq.channel = ch

	q, err := ch.QueueDeclare(
		isq.config.Name,
		isq.config.Durable,
		isq.config.AutoDelete,
		isq.config.Exclusive,
		isq.config.NoWait,
		isq.config.Args,
	)
	if err != nil {
		isq.logger.Fatalf("Failed to declare a queue: %s", err)
		panic(err)
	}
	isq.queue = &q
}

func (isq *ImageSetQueue) Close() {
	isq.conn.Close()
	isq.channel.Close()
}

func (isq *ImageSetQueue) Consume() <-chan amqp.Delivery {
	msgs, err := isq.channel.Consume(
		isq.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		isq.logger.Fatalf("Failed to consume messages from queue: %s", err)
		panic(err)
	}

	return msgs
}

func (isq *ImageSetQueue) Publish(body []byte) error {
	err := isq.channel.Publish(
		"",
		isq.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		isq.logger.Printf("Failed to publish message: %s", err)
		return err
	}
	return nil
}
