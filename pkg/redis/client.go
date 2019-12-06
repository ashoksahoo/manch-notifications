package redis

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var url = os.Getenv("REDIS_URL")

var Pool *redis.Pool

func formatKey(key string) string {
	var env string
	env = os.Getenv("env")
	if env == "" {
		env = "development"
	}
	return "manch-notification:" + env + ":" + key
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	if url == "" {
		url = "redis://redis/"
	}
	Pool = newPool(url)
	Ping()
}
