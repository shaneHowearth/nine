package rabbit

import (
	"fmt"
	"log"
	"strconv"
	"time"

	readarticle "github.com/shanehowearth/nine/readarticles/internal/readarticleservice"
	"github.com/streadway/amqp"
)

// MQ -
type MQ struct {
	conn *amqp.Connection
}

// Connect -
func (r *MQ) Connect() (err error) {
	// Infinite loop.
	// This code tries to connect to the specified MQ server
	// 5 times with a 1 second break between each try.
	// If all of those attempts fail sleep for 5 seconds before starting another round of attempts
	for {
		for i := 0; i < 5; i++ {
			//TODO Get this from env
			r.conn, err = amqp.Dial("amqp://guest:guest@rabbitmq-server:5672/")
			if err == nil {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
		log.Printf("Unable to contact RabbitMQ server after 5 tries, will sleep for 5 seconds before trying again.")
		log.Printf("RabbitMQ connection error: %v", err)
		time.Sleep(5 * time.Second)
	}
}

// Receive -
func (r *MQ) Receive(a *readarticle.Server) error {
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
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer with error: %v", err)
	}
	forever := make(chan bool)

	go func() {
		log.Printf("Waiting for messages")
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			id, err := strconv.Atoi(string(d.Body))
			if err != nil {

				log.Printf("Error fetching new article in MQ Receive: id: %s, err: %v", d.Body, err)
				continue
			}
			article, err := a.Store.FetchOne(id)
			if err != nil {
				// log the error
				log.Printf("Error fetching new article in MQ Receive: id: %s, err: %v", d.Body, err)
				// move on to the next one
				continue
			}
			err = a.Cache.Populate(article)
			if err != nil {
				// log the error
				log.Printf("Error populating new article in MQ Receive: id: %s, err: %v", d.Body, err)
				// move on to the next one
			}
		}
	}()

	log.Printf("Receive start success")
	<-forever
	return nil
}
