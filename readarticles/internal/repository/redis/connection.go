package rediscache

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gomodule/redigo/redis"
)

// Redis -
type Redis struct {
	Pool *redis.Pool
}

var once sync.Once

func (r *Redis) initPool() {
	// Use sync.Once to ensure that only one instance of the pool is ever created.
	// Aka a Singleton.
	const server = "localhost"
	const port = "6379"
	once.Do(func() {
		r.Pool = &redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", fmt.Sprintf(server+":"+port))
				if err != nil {
					log.Printf("ERROR: fail init redis pool: %s", err.Error())
					os.Exit(1)
				}
				return conn, err
			},
		}
	})
}

func (r *Redis) ping(conn redis.Conn) {
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		log.Fatalf("ERROR: failed to ping redis connection: %s", err.Error())
	}
}
