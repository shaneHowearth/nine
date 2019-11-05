package rediscache

import (
	"fmt"
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
	grpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
)

type tmpStruct struct {
	ID    string `redis:"id"`
	Title string `redis:"title"`
	Date  string `redis:"date"`
	Body  string `redis:"body"`
}

// Populate - Fill the cache with data
func (r *Redis) Populate(articles ...*grpcProto.Article) error {
	// get conn and put back when exit from method
	var conn redis.Conn
	if r.Pool == nil {
		r.initPool()
		conn = r.Pool.Get()
		r.ping(conn)
	}
	if conn == nil {
		conn = r.Pool.Get()
	}
	defer conn.Close()

	for i := range articles {
		d1 := strings.Split(articles[i].Date, "T")[0]
		_, err := conn.Do("HSET", articles[i].Id, "title", articles[i].Title, "date", d1, "body", articles[i].Body)
		if err != nil {
			return fmt.Errorf("unable to insert %v with error %v", articles[i], err)
		}
		for _, tag := range articles[i].Tags {

			_, err := conn.Do("RPUSH", tag+":"+d1, articles[i].Id)
			if err != nil {
				return fmt.Errorf("unable to insert %v with error %v", tag, err)
			}
			_, err = conn.Do("RPUSH", "list:"+articles[i].Id, tag)
			if err != nil {
				return fmt.Errorf("unable to insert %v with error %v", tag, err)
			}
		}
	}
	return nil
}

// GetTagInfo -
func (r *Redis) GetTagInfo(tagName, date string) *grpcProto.TagInfo {

	var conn redis.Conn
	if r.Pool == nil {
		r.initPool()
		conn = r.Pool.Get()
		r.ping(conn)
	}
	if conn == nil {
		conn = r.Pool.Get()
	}
	defer conn.Close()
	ids, err := redis.Strings(conn.Do("LRANGE", tagName+":"+date, 0, 2147483647))
	if err != nil {
		// TODO
		// return nil, fmt.Errorf("unable to insert %v with error %v", tagName+":"+date, err)
		log.Println("Error me")
	}
	tags := make(map[string]struct{})
	for _, id := range ids {
		taglist, err := r.GetTags(id)
		if err != nil {
			// TODO
			// return nil, fmt.Errorf("unable to insert %v with error %v", id, err)
			log.Println("Error me")
		}
		for _, tag := range taglist {
			tags[tag] = struct{}{}
		}
	}
	// TODO
	return nil
}

// GetTags -
func (r *Redis) GetTags(id string) ([]string, error) {
	var conn redis.Conn
	if r.Pool == nil {
		r.initPool()
		conn = r.Pool.Get()
		r.ping(conn)
	}
	if conn == nil {
		conn = r.Pool.Get()
	}
	defer conn.Close()
	return redis.Strings(conn.Do("LRANGE", "list:"+id, 0, 2147483647))
}

// GetByID -
func (r *Redis) GetByID(id string) *grpcProto.Article {
	// get conn and put back when exit from method
	var conn redis.Conn
	if r.Pool == nil {
		r.initPool()
		conn = r.Pool.Get()
		r.ping(conn)
	}
	if conn == nil {
		conn = r.Pool.Get()
	}
	defer conn.Close()

	dataset, err := redis.Values(conn.Do("HGETALL", id))
	if err != nil {
		log.Printf("ERROR: failed get key %s, error %s", id, err.Error())
		return &grpcProto.Article{}
	}

	// Put dataset into an Article
	a := grpcProto.Article{}
	f := tmpStruct{}

	if len(dataset) == 0 {
		log.Printf("Cache miss looking for %s", id)
		// Cache miss
		// Check DB in case it exists
		return &a
	}
	err = redis.ScanStruct(dataset, &f)
	if err != nil {
		log.Printf("error scanning struct: %v", err)
	}
	a.Id = id
	a.Title = f.Title
	a.Date = f.Date
	a.Body = f.Body
	a.Tags, err = r.GetTags(id)
	if err != nil {
		log.Printf("Shit %v, %T", err, []byte(dataset[3].([]byte)))
	}
	return &a
}
