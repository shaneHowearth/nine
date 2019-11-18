package rabbit

import "github.com/streadway/amqp"

// SetAMQPDialForTest - Allow test functions to set the Rabbit MQ dialer function
func SetAMQPDialForTest(d func(string) (*amqp.Connection, error)) {
	amqpDial = d
}
