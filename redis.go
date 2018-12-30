package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func main() {
	var pool = newPool()
	c := pool.Get()
	n1, _ := c.Do("SET", "key", []byte("value"))
	n2, _ := redis.Bytes(c.Do("GET", "key"))
	fmt.Println(n1, string(n2))
	defer c.Close()
}
