package rediscache_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
	rediscache "github.com/shanehowearth/nine/readarticles/internal/repository/redis"
	"github.com/stretchr/testify/assert"
)

// Create global (to these tests) Redis and connection objects
var (
	testR rediscache.Redis
	conn  *redigomock.Conn
)

func TestMain(m *testing.M) {

	conn = redigomock.NewConn()
	mockPool := &redis.Pool{
		Dial:    func() (redis.Conn, error) { return conn, nil },
		MaxIdle: 10,
	}
	testR = rediscache.Redis{Pool: mockPool}
	os.Exit(m.Run())
}

func TestGetByID(t *testing.T) {
	testcases := map[string]struct {
		title string
	}{
		"Error from redis": {title: "1"},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			conn.Clear()
			conn.Command("HGETALL").ExpectError(fmt.Errorf("unexpected Error"))

			output, found := testR.GetByID(tc.title)

			assert.Equal(t, &grpcProto.Article{}, output, "Article returned did not match expected")
			assert.True(t, found, "Expected to find article")

		})
	}
}
