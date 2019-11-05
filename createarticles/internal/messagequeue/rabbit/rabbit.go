package rabbit

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// MQ -
type MQ struct {
	conn *amqp.Connection
}

// Connect -
func (r *MQ) Connect() (err error) {
	//TODO Get this from env
	for {
		for i := 0; i < 5; i++ {
			r.conn, err = amqp.Dial("amqp://guest:guest@rabbitmq-server:5672/")
			if err == nil {
				return
			}
		}
		log.Printf("Trouble connecting to RabbitMQ, error: %v", err)
		time.Sleep(5 * time.Second)
	}
}

// Publish -
func (r *MQ) Publish(id string) error {
	if r.conn == nil {
		err := r.Connect()
		if err != nil {
			return fmt.Errorf("failed to create connection with error: %v", err)
		}
	}
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel with error: %v", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"articles", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue with error: %v", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(id),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message with error: %v", err)

	}
	return nil
}
