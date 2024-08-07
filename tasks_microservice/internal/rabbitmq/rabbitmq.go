package rabbitmq

import (
	"fmt"
	"log"
	"tasks_microservice/internal/handlers"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn *amqp.Connection
}

// ConnectRabbitMQ creates a new RabbitMQ connection.
func ConnectRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return &RabbitMQ{Conn: conn}, nil
}

// SetupRabbitMQConsumers sets up RabbitMQ consumers for the given queues.
func (r *RabbitMQ) SetupRabbitMQConsumers(handler *handlers.Handlers) {
	r.declareQueue("tasks_queue")
	go r.consume("tasks_queue", handler.HandleTaskValues)
}

// Publish sends a message to the specified RabbitMQ queue.
func (r *RabbitMQ) Publish(queueName string, body []byte, replyTo string, correlationId string) error {
	channel, err := r.Conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer channel.Close()

	err = channel.Publish(
		"",          // exchange
		queueName,   // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			ReplyTo:       replyTo,
			CorrelationId: correlationId,
		},
	)
	if err != nil {
		log.Printf("Error publishing message to RabbitMQ: %v", err)
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}

// declareQueue declares a RabbitMQ queue.
func (r *RabbitMQ) declareQueue(queueName string) {
	channel, err := r.Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		queueName,  // queue name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
}

// consume starts consuming messages from the specified RabbitMQ queue.
func (r *RabbitMQ) consume(queueName string, handler func(amqp.Delivery)) {
	channel, err := r.Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	msgs, err := channel.Consume(
		queueName,  // queue name
		"",         // consumer tag
		true,       // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()
}

// Close closes the RabbitMQ connection.
func (r *RabbitMQ) Close() error {
	if r.Conn != nil {
		return r.Conn.Close()
	}
	return nil
}
