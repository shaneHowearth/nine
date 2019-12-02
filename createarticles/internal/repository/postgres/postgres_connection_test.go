package datastore_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/bouk/monkey"
	datastore "github.com/shanehowearth/nine/createarticles/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	testcases := map[string]struct {
		retry     int
		uri       string
		err       error
		willPanic bool
	}{
		"Happy Path":   {retry: 1, uri: "postgres://pqgotest:password@localhost/pqgotest?sslmode=disable"},
		"No Retry set": {uri: "postgres://pqgotest:password@localhost/pqgotest?sslmode=disable"},
		"No URI set":   {retry: 1, willPanic: true, err: fmt.Errorf("no Postgres URI configured")},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			p := datastore.Postgres{Retry: tc.retry, URI: tc.uri}
			datastore.SetSQLOpenForTest(func(dbtype, dbURI string) (*sql.DB, error) {
				return nil, tc.err
			})
			if tc.willPanic {
				fakeLogPanic := func(msg ...interface{}) {
					assert.Equal(t, tc.err.Error(), msg[0])
					panic(tc.err.Error())
				}
				patch := monkey.Patch(log.Panic, fakeLogPanic)
				defer patch.Unpatch()
				assert.PanicsWithValue(t, tc.err.Error(), func() { _ = p.Connect() }, "log.Panic was not called")
			} else {
				output := p.Connect()
				assert.Nil(t, output, "Got an unexpected error from connect, %v", output)
			}
		})
	}
}
