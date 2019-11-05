package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	grpcProto "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"

	"github.com/lib/pq"
)

// Postgres -
type Postgres struct {
	db *sql.DB
}

// Connect - Create the connection to the database
func (p *Postgres) Connect() (err error) {
	connStr := "postgres://articlewriter:hu8jmn3@articledb/article_postgres_db?sslmode=disable"
	p.db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to open connection to postgres, error: %v", err)
	}

	return nil
}

// CreateOrUpdate -
func (p *Postgres) Create(article *grpcProto.Article) (string, error) {
	if p.db == nil {
		perr := p.Connect()
		if perr != nil {
			log.Fatalf("unable to connect to postgres server with error: %v", perr)
		}
	}
	var id int
	err := p.db.QueryRow(`INSERT INTO article(title, pub_date, body, tags)
	VALUES($1, $2, $3, $4) RETURNING id`, article.GetTitle(), article.GetDate(), article.GetBody(), pq.Array(article.GetTags())).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("postgres CreateOrUpdate returned error: %v", err)
	}
	return strconv.Itoa(id), nil
}
