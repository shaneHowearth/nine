package articleservice

import (
	"context"
	"log"

	grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
	repo "github.com/shanehowearth/nine/readarticles/integration/repository/cache/v1"
	database "github.com/shanehowearth/nine/readarticles/integration/repository/database/v1"
)

// articleServiceServer is implementation of grpcProto.ArticleServiceServer proto interface
type articleServiceServer struct {
	Cache repo.Cache
	Store database.Storage
}

// NewArticleService creates Article service
func NewArticleService(c repo.Cache, s database.Storage) grpcProto.ArticleServiceServer {
	if c == nil {
		log.Panic("Cache supplied for NewArticleService is nil")
	}
	if s == nil {
		log.Panic("Store supplied for NewArticleService is nil")
	}
	a := articleServiceServer{Cache: c, Store: s}
	// Fill the cache with data
	const startcount = 100
	articles, err := a.Store.FetchLatestRows(startcount)
	if err != nil {
		log.Panicf("Unable to get any data for the cache, error: %v", err)
	}
	if err = a.Cache.Populate(articles...); err != nil {
		log.Panicf("Unable to populate cache, error: %v", err)
	}
	return &a
}

// GetTagInfo -
func (a *articleServiceServer) GetTagInfo(ctx context.Context, req *grpcProto.ArticleRequest) (*grpcProto.TagInfo, error) {
	log.Printf("GetArticle readarticleservice req: %#+v", req)
	article := a.Cache.GetTagInfo(req.GetTag(), req.GetDate())

	return article, nil
}

// GetArticle -
func (a *articleServiceServer) GetArticle(ctx context.Context, req *grpcProto.ArticleRequest) (*grpcProto.Article, error) {
	log.Printf("GetArticle readarticleservice req: %#+v", req)
	article := a.Cache.GetByID(req.GetId())

	return article, nil
}
