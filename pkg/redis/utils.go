package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func AcquireLock(key, value string, ttl int) (interface{}, error) {
	key = formatKey(key)
	conn := Pool.Get()
	defer conn.Close()
	result, err := conn.Do("SETNX", key, value)
	if result == 1 {
		_, err := conn.Do("EXPIRE", key, ttl)
		if err != nil {
			fmt.Println("Error to set expire")
		}
	}
	return result, err
}

func ReleaseLock(key string) error {
	key = formatKey(key)
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func Get(key string) (interface{}, error) {
	key = formatKey(key)
	conn := Pool.Get()
	defer conn.Close()
	return conn.Do("GET", key)
}

func Ping() error {

	conn := Pool.Get()
	defer conn.Close()

	result, err := redis.String(conn.Do("PING"))
	fmt.Println("PING REDIS: ", result)
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

func Set(key string, value []byte) error {
	key = formatKey(key)
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}
