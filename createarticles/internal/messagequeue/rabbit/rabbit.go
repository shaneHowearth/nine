package rabbit

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

// MQ -
type MQ struct {
	conn  *amqp.Connection
	Retry int
	URI   string
}

var amqpDial = amqp.Dial

// Connect -
func (mq *MQ) Connect() (err error) {
	// Retry MUST be > 0
	if mq.Retry == 0 {
		log.Printf("Cannot use a Retry of zero, this process will to default retry to 5")
		mq.Retry = 5
	}

	// Note: Even though amqp.ParseURI(uri) will validate the URI formed, check here that the minimum required exists
	if mq.URI == "" {
		log.Printf("No Message Queue URI configured")
	}

	for {
		for i := 0; i < mq.Retry; i++ {
			mq.conn, err = amqpDial(mq.URI)
			if err == nil {
				// Successful connection
				return nil
			}
			time.Sleep(1 * time.Second)
		}
		log.Printf("Trouble connecting to RabbitMQ, error: %v", err)
		time.Sleep(5 * time.Second)
	}
}

// Publish -
func (mq *MQ) Publish(id string) error {
	if mq.conn == nil {
		err := mq.Connect()
		if err != nil {
			return fmt.Errorf("failed to create connection with error: %v", err)
		}
	}
	ch, err := mq.conn.Channel()
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
