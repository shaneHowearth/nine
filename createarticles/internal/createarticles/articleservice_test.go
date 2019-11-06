package articleservice_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/bouk/monkey"
	grpcProto "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"
	messagequeue "github.com/shanehowearth/nine/createarticles/integration/messagequeue/v1"
	repository "github.com/shanehowearth/nine/createarticles/integration/repository/v1"
	SUT "github.com/shanehowearth/nine/createarticles/internal/createarticles"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct{}

func (m *mockRepo) Create(article *grpcProto.Article) (s string, e error) {
	return storeString, storeError
}

type mockmessageQueue struct{}

var storeError error
var storeString string
var signalError error

func (m2 *mockmessageQueue) Publish(id string) error { return signalError }

func TestNewArticleService(t *testing.T) {
	mockStore := &mockRepo{}
	mockMQ := &mockmessageQueue{}
	testcases := map[string]struct {
		store       repository.Storage
		mq          messagequeue.MQ
		server      SUT.Server
		errMessage  string
		expectPanic bool
	}{
		"Happy Path": {store: mockStore,
			mq:     mockMQ,
			server: SUT.Server{Storage: mockStore, Signal: mockMQ},
		},
		"Missing Storage": {mq: mockMQ, errMessage: "NewArticleService has no cache to get articles from", expectPanic: true},
		"Missing MQ":      {store: mockStore, errMessage: "NewArticleService has no messagequeue to send to", expectPanic: true},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			if tc.expectPanic {
				fakeLogFatal := func(msg ...interface{}) {
					assert.Equal(t, tc.errMessage, msg[0])
					panic("log.Fatal called")
				}
				patch := monkey.Patch(log.Fatal, fakeLogFatal)
				defer patch.Unpatch()
				assert.PanicsWithValue(t, "log.Fatal called", func() { SUT.NewArticleService(tc.store, tc.mq) }, "log.Fatal was not called")
			} else {
				output := SUT.NewArticleService(tc.store, tc.mq)
				assert.Equal(t, *output, tc.server)
			}
		})
	}
}

func TestCreateArticle(t *testing.T) {
	mockStore := &mockRepo{}
	mockMQ := &mockmessageQueue{}
	testcases := map[string]struct {
		ctx        context.Context
		input      *grpcProto.Article
		response   grpcProto.Acknowledgement
		storeID    string
		errMessage string
		signalErr  error
		storeErr   error
	}{
		"Happy Path": {ctx: context.Background(),
			storeID:  "1",
			input:    &grpcProto.Article{},
			response: grpcProto.Acknowledgement{Id: "1"},
		},

		"Signal Problem":  {errMessage: "", storeID: "1", signalErr: fmt.Errorf("random error message"), response: grpcProto.Acknowledgement{Id: "1"}},
		"Storage Problem": {errMessage: "", storeID: "1", storeErr: fmt.Errorf("random error message"), response: grpcProto.Acknowledgement{}},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ss := SUT.NewArticleService(mockStore, mockMQ)
			signalError = tc.signalErr
			storeError = tc.storeErr
			storeString = tc.storeID
			output, err := ss.CreateArticle(tc.ctx, tc.input)
			if tc.storeErr == nil {
				assert.Nil(t, err, "Was not expecting an error got %v", err)
				assert.Equal(t, tc.storeID, (*output).Id, "Expecting Acknowledgement id to be the same")
			} else {
				assert.NotNil(t, err, "Was expecting an error %v", err)
				assert.NotEqual(t, tc.storeID, (*output).Id, "Expecting Acknowledgement id to be the same")
			}
			assert.Equal(t, tc.response, *output, "Expecting Acknowledgement")
		})
	}

}
