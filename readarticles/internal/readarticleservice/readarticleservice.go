package articleservice

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
	repo "github.com/shanehowearth/nine/readarticles/integration/repository/cache/v1"
	database "github.com/shanehowearth/nine/readarticles/integration/repository/database/v1"
)

// Server -
type Server struct {
	Cache repo.Cache
	Store database.Storage
}

// NewArticleService -
func NewArticleService(c repo.Cache, s database.Storage) *Server {
	if c == nil {
		log.Fatal("Cache supplied for NewArticleService is nil")
	}
	if s == nil {
		log.Fatal("Store supplied for NewArticleService is nil")
	}
	a := Server{Cache: c, Store: s}
	// Fill the cache with data
	const startcount = 100
	var err error
	var articles []*grpcProto.Article
	for numTries := 0; numTries < 5; numTries++ {
		articles, err = a.Store.FetchLatestRows(startcount)
		if err == nil {
			// No error means we have succeeded in talking to the DB
			break
		}
		// Wait a second for the DB to come back to life
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Printf("Unable to get any data for the cache, error: %v", err)
	}
	if err = a.Cache.Populate(articles...); err != nil {
		log.Printf("Unable to populate cache, error: %v", err)
	}
	return &a
}

// GetTagInfo -
func (a *Server) GetTagInfo(ctx context.Context, req *grpcProto.ArticleRequest) (*grpcProto.TagInfo, error) {
	tagInfo := a.Cache.GetTagInfo(req.GetTag(), req.GetDate())

	return tagInfo, nil
}

// GetArticle -
func (a *Server) GetArticle(ctx context.Context, req *grpcProto.ArticleRequest) (*grpcProto.Article, error) {
	id := req.GetId()
	article, found := a.Cache.GetByID(id)
	if !found {
		iid, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Bad id supplied %s", id)
			return &grpcProto.Article{}, fmt.Errorf("bad id supplied %q", id)
		}
		article, _ = a.Store.FetchOne(iid)
		if article == nil {
			log.Println("Got nil back from FetchOne")
			return &grpcProto.Article{}, nil
		}
		if err = a.Cache.Populate(article); err != nil {
			log.Printf("Unable to populate cache with error: %v", err)
		}
	}
	return article, nil
}
