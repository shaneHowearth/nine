package createclient

import (
	"context"
	"log"
	"time"

	proto "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"
	"google.golang.org/grpc"
)

// CreateClient -
type CreateClient struct {
	Address string
}

func (s *CreateClient) newConnection() (proto.ArticleServiceClient, *grpc.ClientConn) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return proto.NewArticleServiceClient(conn), conn
}

// CreateArticle -
func (s *CreateClient) CreateArticle(article *proto.Article) *proto.Acknowledgement {
	c, conn := s.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ack, err := c.CreateArticle(ctx, article)
	if err != nil {
		// TODO: Log error
		log.Printf("Got %v", err)
	}
	return ack
}
