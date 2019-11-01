package repository

// Provide interfaces for the repository cache to use in order for the
// microservice to communicate with the cache.
// Use interfaces here so the cache depends on the microservice, not the other way round
// allowing any cache to be dropped in, as long as it satisfies the interface(s)

import (
	grpcProto "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"
)

type Storage interface {
	CreateorUpdate(article *grpcProto.Article) (string, error)
}
