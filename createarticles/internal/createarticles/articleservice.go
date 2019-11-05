package articleservice

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	grpc "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"
	messagequeue "github.com/shanehowearth/nine/createarticles/integration/messagequeue/v1"
	repository "github.com/shanehowearth/nine/createarticles/integration/repository/v1"
)

// articleServiceServer is implementation of v1.ArticleServiceServer proto interface.
type articleServiceServer struct {
	Storage repository.Storage
	Signal  messagequeue.MQ
}

// NewArticleService creates Article service.
func NewArticleService(s repository.Storage, mq messagequeue.MQ) grpc.ArticleServiceServer {
	if s == nil {
		log.Panic("NewArticleService has no cache to get articles from")
	}
	if mq == nil {
		log.Panic("NewArticleService has no messagequeue to send to")

	}
	a := articleServiceServer{Storage: s, Signal: mq}
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
	// Alert all the observers that a new article exists
	if err = a.Signal.Publish(id); err != nil {
		log.Printf("Error sending alert of new article: %v", err)
	}
	return &grpc.Acknowledgement{Id: id}, nil
}
