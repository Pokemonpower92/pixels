package queue

import amqp "github.com/rabbitmq/amqp091-go"

type Queue interface {
	Initialize()
	Close()
	Consume() <-chan amqp.Delivery
	Publish([]byte) error
}
