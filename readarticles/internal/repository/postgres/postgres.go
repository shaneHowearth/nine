package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	// Infinite loop
	for {
		for i := 0; i < 5; i++ {
			p.db, err = sql.Open("postgres", connStr)
			if err == nil {
				return
			}
		}
		// TODO raise monitored alert
		log.Printf("Unable to open connection to postgres, error: %v", err)
		time.Sleep(5 * time.Second)
	}
}

// FetchLatestRows -
func (p *Postgres) FetchLatestRows(n int) ([]*grpcProto.Article, error) {
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

	return p.toArticle(rows), nil
}

func (p *Postgres) toArticle(rows *sql.Rows) (articles []*grpcProto.Article) {
	var id, title, date, body sql.NullString
	var tags []sql.NullString
	for rows.Next() {
		err := rows.Scan(&id, &title, &date, &body, pq.Array(&tags))
		if err != nil {
			log.Printf("scan error in toArticle: %v", err)
			continue
		}
		var tagStrings []string
		for _, t := range tags {
			tagStrings = append(tagStrings, t.String)
		}
		articles = append(articles, &grpcProto.Article{Id: id.String, Title: title.String, Date: date.String, Body: body.String, Tags: tagStrings})
	}
	return articles
}

// FetchOne -
func (p *Postgres) FetchOne(aid int) (article *grpcProto.Article, err error) {
	if p.db == nil {
		perr := p.Connect()
		if perr != nil {
			log.Fatalf("unable to connect to postgres server with error: %v", perr)
		}
	}

	rows, err := p.db.Query(`SELECT id, title, pub_date, body, tags FROM article WHERE id = $1;`, aid)
	if err != nil {
		return nil, fmt.Errorf("postgres FetchOne returned error: %v", err)
	}
	articles := p.toArticle(rows)
	if len(articles) == 1 {
		return articles[0], nil
	}
	return
}
