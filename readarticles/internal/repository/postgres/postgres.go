package datastore

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
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

// FetchLatestRows -
func (p *Postgres) FetchLatestRows(n int) (articles []*grpcProto.Article, err error) {
	if p.db == nil {
		perr := p.Connect()
		if perr != nil {
			log.Fatalf("unable to connect to postgres server with error: %v", perr)
		}
	}
	rows, err := p.db.Query(`SELECT id, title, pub_date, body, tags FROM article ORDER BY id DESC LIMIT $1;`, n)
	if err != nil {
		return nil, fmt.Errorf("postgres FetchLatestRows returned error: %v", err)
	}

	var id, title, date, body sql.NullString
	var tags []sql.NullString
	for rows.Next() {
		log.Printf("ROWS: %#+v", rows)
		if err := rows.Scan(&id, &title, &date, &body, pq.Array(&tags)); err != nil {
			return nil, fmt.Errorf("failed to scan with error: %v", err)
		}
		var tagStrings []string
		for _, t := range tags {
			tagStrings = append(tagStrings, t.String)
		}
		articles = append(articles, &grpcProto.Article{Id: id.String, Title: title.String, Date: date.String, Body: body.String, Tags: tagStrings})
	}
	return articles, nil
}

// FetchOne -
func (p *Postgres) FetchOne(aid int) (*grpcProto.Article, error) {
	if p.db == nil {
		perr := p.Connect()
		if perr != nil {
			log.Fatalf("unable to connect to postgres server with error: %v", perr)
		}
	}
	var id, title, date, body string
	var tags []string
	row, err := p.db.Query(`SELECT id, title, pub_date, body, tags FROM article WHERE id = $1;`, aid)
	if err != nil {
		return nil, fmt.Errorf("postgres FetchOne returned error: %v", err)
	}
	if err = row.Scan(&id, &title, &date, &body, &tags); err != nil {
		return nil, fmt.Errorf("scan in FetchOne returned error: %v", err)

	}

	return &grpcProto.Article{Id: id, Title: title, Date: date, Body: body, Tags: tags}, nil
}
