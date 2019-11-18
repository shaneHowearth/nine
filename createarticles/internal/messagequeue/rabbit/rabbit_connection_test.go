package rabbit_test

import (
	"testing"

	"github.com/shanehowearth/nine/createarticles/internal/messagequeue/rabbit"
	"github.com/stretchr/testify/assert"

	"github.com/streadway/amqp"
)

func TestConnect(t *testing.T) {

	testcases := map[string]struct {
		retry int
		uri   string
		err   error
	}{
		"Happy Path":    {retry: 1, uri: "amqp://localhost:5672/%2f"},
		"Missing Retry": {uri: "amqp://localhost:5672/%2f"},
		"Missing URI":   {retry: 1},
		// In order to test further the infinite loop needs to be stopped
		// "Bad connection": {retry: 1, uri: "amqp://localhost:5672/%2f", err: fmt.Errorf("bad connection")},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			rabbit.SetAMQPDialForTest(func(msg string) (*amqp.Connection, error) {
				return nil, tc.err
			})
			mq := rabbit.MQ{Retry: tc.retry, URI: tc.uri}
			output := mq.Connect()
			assert.Nil(t, output, "Got an unexpected error from connect, %v", output)
		})
	}
}
