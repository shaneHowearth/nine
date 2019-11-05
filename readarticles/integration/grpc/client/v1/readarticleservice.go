package readarticleclient

import (
	"context"
	"log"
	"time"

	grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
	"google.golang.org/grpc"
)

// ReadArticleClient -
type ReadArticleClient struct {
	Address string
}

func (s *ReadArticleClient) newConnection() (grpcProto.ArticleServiceClient, *grpc.ClientConn) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return grpcProto.NewArticleServiceClient(conn), conn
}

// GetArticle -
func (s *ReadArticleClient) GetArticle(ar *grpcProto.ArticleRequest) (*grpcProto.Article, error) {
	c, conn := s.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.GetArticle(ctx, ar)
}

// GetTagInfo
func (s *ReadArticleClient) GetTagInfo(ar *grpcProto.ArticleRequest) (*grpcProto.TagInfo, error) {
	c, conn := s.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.GetTagInfo(ctx, ar)
}
