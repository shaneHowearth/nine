package main

import (
	"context"
	"log"
	"net"
	"os"

	article "github.com/shanehowearth/nine/createarticles/internal/createarticles"
	repo "github.com/shanehowearth/nine/createarticles/internal/repository/postgres"

	grpcProto "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

var ss = article.NewArticleService(new(repo.Postgres))

func main() {

	portNum := os.Getenv("PORT_NUM")
	lis, err := net.Listen("tcp", "0.0.0.0:"+portNum)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcProto.RegisterArticleServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// GetArticle -
func (s *server) CreateArticle(ctx context.Context, req *grpcProto.Article) (*grpcProto.Acknowledgement, error) {
	st, err := ss.CreateArticle(ctx, req)
	if err != nil {
		return nil, err
	}
	return st, nil
}
