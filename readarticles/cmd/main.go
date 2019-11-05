package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/shanehowearth/nine/readarticles/internal/messagequeue/rabbit"
	readarticle "github.com/shanehowearth/nine/readarticles/internal/readarticleservice"
	database "github.com/shanehowearth/nine/readarticles/internal/repository/postgres"
	repo "github.com/shanehowearth/nine/readarticles/internal/repository/redis"

	grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

var ss = readarticle.NewArticleService(new(repo.Redis), new(database.Postgres))

func main() {

	portNum := os.Getenv("PORT_NUM")
	lis, err := net.Listen("tcp", "0.0.0.0:"+portNum)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	mq := rabbit.MQ{}
	go func() {
		if err := mq.Receive(ss); err != nil {
			log.Printf("mq Receiver died with error: %v", err)
		}
	}()
	s := grpc.NewServer()
	grpcProto.RegisterArticleServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// GetTagInfo -
func (s *server) GetTagInfo(ctx context.Context, req *grpcProto.ArticleRequest) (*grpcProto.TagInfo, error) {
	st, err := ss.GetTagInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return st, nil
}

// GetArticle -
func (s *server) GetArticle(ctx context.Context, req *grpcProto.ArticleRequest) (*grpcProto.Article, error) {
	st, err := ss.GetArticle(ctx, req)
	if err != nil {
		return nil, err
	}
	return st, nil
}
