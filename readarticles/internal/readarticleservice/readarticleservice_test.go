package articleservice_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/bouk/monkey"
	grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
	repo "github.com/shanehowearth/nine/readarticles/integration/repository/cache/v1"
	database "github.com/shanehowearth/nine/readarticles/integration/repository/database/v1"
	SUT "github.com/shanehowearth/nine/readarticles/internal/readarticleservice"
	"github.com/stretchr/testify/assert"
)

type mockRepoCache struct{}

var cacheArticle *grpcProto.Article
var cacheFound bool

func (m *mockRepoCache) GetByID(id string) (*grpcProto.Article, bool) { return cacheArticle, cacheFound }

var tagInfo *grpcProto.TagInfo

func (m *mockRepoCache) GetTagInfo(tag, date string) *grpcProto.TagInfo { return tagInfo }

var populateErr error

func (m *mockRepoCache) Populate(...*grpcProto.Article) error { return populateErr }

type mockStorage struct{}

var fetchLatestRowsErr error

func (ms *mockStorage) FetchLatestRows(n int) (as []*grpcProto.Article, e error) {
	return as, fetchLatestRowsErr
}

var fetchOneError error
var fetchOneArticle *grpcProto.Article

func (ms *mockStorage) FetchOne(id int) (a *grpcProto.Article, e error) {
	return fetchOneArticle, fetchOneError
}

func TestNewArticleService(t *testing.T) {
	mockStore := &mockStorage{}
	mockCache := &mockRepoCache{}
	testcases := map[string]struct {
		cache       repo.Cache
		store       database.Storage
		server      SUT.Server
		errMessage  string
		expectPanic bool
		storeErr    error
		cacheErr    error
	}{
		"Happy Path":          {cache: mockCache, store: mockStore, server: SUT.Server{Cache: mockCache, Store: mockStore}},
		"Missing Cache":       {store: mockStore, expectPanic: true, errMessage: "Cache supplied for NewArticleService is nil"},
		"Missing Store":       {cache: mockCache, expectPanic: true, errMessage: "Store supplied for NewArticleService is nil"},
		"Store returns error": {cache: mockCache, store: mockStore, server: SUT.Server{Cache: mockCache, Store: mockStore}, storeErr: fmt.Errorf("error returned")},
		"Cache returns error": {cache: mockCache, store: mockStore, server: SUT.Server{Cache: mockCache, Store: mockStore}, cacheErr: fmt.Errorf("error returned")},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			fetchLatestRowsErr = tc.storeErr
			populateErr = tc.cacheErr
			if tc.expectPanic {
				fakeLogFatal := func(msg ...interface{}) {
					assert.Equal(t, tc.errMessage, msg[0])
					panic("log.Fatal called")
				}
				patch := monkey.Patch(log.Fatal, fakeLogFatal)
				defer patch.Unpatch()
				assert.PanicsWithValue(t, "log.Fatal called", func() { SUT.NewArticleService(tc.cache, tc.store) }, "log.Fatal was not called")
			} else {

				output := SUT.NewArticleService(tc.cache, tc.store)
				assert.Equal(t, *output, tc.server)
			}
		})
	}
}

func TestTagInfo(t *testing.T) {
	mockStore := &mockStorage{}
	mockCache := &mockRepoCache{}

	testcases := map[string]struct {
		ctx      context.Context
		input    *grpcProto.ArticleRequest
		response *grpcProto.TagInfo
	}{
		"Happy Path": {
			ctx:      context.Background(),
			input:    &grpcProto.ArticleRequest{},
			response: &grpcProto.TagInfo{},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			tagInfo = tc.response
			ss := SUT.NewArticleService(mockCache, mockStore)
			output, err := ss.GetTagInfo(tc.ctx, tc.input)
			assert.Equal(t, tc.response, output, "Expected %v got %v", tc.response, output)
			assert.Nil(t, err, "Not expecting an error")
		})
	}
}

func TestGetArticle(t *testing.T) {
	mockStore := &mockStorage{}
	mockCache := &mockRepoCache{}

	testcases := map[string]struct {
		ctx            context.Context
		input          *grpcProto.ArticleRequest
		response       *grpcProto.Article
		fetchedArticle *grpcProto.Article
		found          bool
		errorReturned  bool
		cacheErr       error
	}{
		"Happy Path": {
			ctx:      context.Background(),
			input:    &grpcProto.ArticleRequest{Id: "1"},
			response: &grpcProto.Article{Id: "1"},
			found:    true,
		},
		"No Article found": {
			ctx:      context.Background(),
			input:    &grpcProto.ArticleRequest{Id: "1"},
			response: &grpcProto.Article{},
		},
		"No Article in Cache, but one in DB": {
			ctx:            context.Background(),
			input:          &grpcProto.ArticleRequest{Id: "1"},
			response:       &grpcProto.Article{Id: "1"},
			fetchedArticle: &grpcProto.Article{Id: "1"},
		},
		"Bad Id supplied": {
			ctx:            context.Background(),
			input:          &grpcProto.ArticleRequest{Id: "Bad"},
			response:       &grpcProto.Article{},
			fetchedArticle: &grpcProto.Article{},
			errorReturned:  true,
		},
		"Unable to populate cache": {
			ctx:            context.Background(),
			input:          &grpcProto.ArticleRequest{Id: "22"},
			response:       &grpcProto.Article{},
			fetchedArticle: &grpcProto.Article{},
			cacheErr:       fmt.Errorf("cache error"),
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {

			ss := SUT.NewArticleService(mockCache, mockStore)
			cacheArticle = tc.response
			cacheFound = tc.found
			fetchOneArticle = tc.fetchedArticle
			populateErr = tc.cacheErr

			output, err := ss.GetArticle(tc.ctx, tc.input)
			assert.Equal(t, *tc.response, *output, "Expected %v got %v", tc.response, output)
			if tc.errorReturned {
				assert.NotNil(t, err, "Expecting error")
			} else {
				assert.Nil(t, err, "Not expecting an error, but got %v", err)
			}
		})
	}
}
