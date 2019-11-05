package articleservice

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	grpc "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"
	repository "github.com/shanehowearth/nine/createarticles/integration/repository/v1"
)

// articleServiceServer is implementation of v1.ArticleServiceServer proto interface.
type articleServiceServer struct {
	Storage repository.Storage
}

// NewArticleService creates Article service.
func NewArticleService(s repository.Storage) grpc.ArticleServiceServer {
	if s == nil {
		log.Panic("NewArticleService has no cache to get articles from")
	}
	a := articleServiceServer{Storage: s}
	return &a
}

// CreateArticle - Create Article.
func (a *articleServiceServer) CreateArticle(ctx context.Context, req *grpc.Article) (*grpc.Acknowledgement, error) {
	// Handler validates the input.
	id, err := a.Storage.Create(req)
	if err != nil {
		// create a unique uuid for the user to quote to tech support.
		id, uuiderr := uuid.NewUUID()
		if uuiderr != nil {
			// This should never happen, but if it does an alert will need to be raised immediately.
			log.Printf("Error creating uuid during article creation with context: %+v, request: %+v. error: %v", ctx, req, uuiderr)
		}
		log.Printf("Error creating article in repository: %v, code: %s", err, id.String())
		return &grpc.Acknowledgement{}, fmt.Errorf("unable to create article, please quote code: %s", id.String())
	}
	return &grpc.Acknowledgement{Id: id}, nil
}
