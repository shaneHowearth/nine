// Provide interface for the database to use for communicating
// with the microservice.

package database

import grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"

// Storage -
type Storage interface {
	FetchLatestRows(n int) ([]*grpcProto.Article, error)
	FetchOne(id int) (*grpcProto.Article, error)
}
