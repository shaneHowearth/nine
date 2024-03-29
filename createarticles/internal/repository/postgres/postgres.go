package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	grpcProto "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"

	"github.com/lib/pq"
)

// Postgres -
type Postgres struct {
	db    *sql.DB
	Retry int
	URI   string
}

var sqlOpen = sql.Open

// Connect - Create the connection to the database
func (p *Postgres) Connect() (err error) {
	// Retry MUST be >= 1
	if p.Retry == 0 {
		log.Print("Cannot use a Retry of zero, this process will to default retry to 5")
		p.Retry = 5
	}
	if p.URI == "" {
		log.Panicf("no Postgres URI configured")
	}

	// Infinite loop
	// Keep trying forever
	for {
		for i := 0; i < p.Retry; i++ {
			p.db, err = sqlOpen("postgres", p.URI)
			if err == nil {
				if pingerr := p.db.Ping(); pingerr != nil {
					log.Printf("Unable to ping database with error %v", pingerr)
				} else {
					// Successful connection
				log.Print("Successfully connected to postgres DB")
					return nil
				}
			}
			time.Sleep(1 * time.Second)
		}

		backoff := time.Duration(p.Retry*rand.Intn(10)) * time.Second
		log.Printf("ALERT: Trouble connecting to Postgres, error: %v, going to re-enter retry loop in %s seconds", err, backoff.String())
		time.Sleep(backoff)
	}
}

// CreateArticle -
func (p *Postgres) CreateArticle(article *grpcProto.Article) (string, error) {
	if p.db == nil {
		perr := p.Connect()
		if perr != nil {
			// should never get here
			log.Fatalf("unable to connect to postgres server with error: %v", perr)
		}
	}
	var id int
	err := p.db.QueryRow(`INSERT INTO article(title, pub_date, body, tags)
	VALUES($1, $2, $3, $4) RETURNING id`, article.GetTitle(), article.GetDate(), article.GetBody(), pq.Array(article.GetTags())).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("postgres CreateArticle returned error: %v", err)
	}
	return strconv.Itoa(id), nil
}
