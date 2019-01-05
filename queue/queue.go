package queue

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Queue represents a proxy queue for requests
type Queue struct {
	Name   string   `json:"name"`
	Offset int      `json:"offset"`
	Jobs   [][]byte `json:"jobs"`
}

type redisPool redis.Pool

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				fmt.Printf("Could not connect to REDIS DB")
				panic(err.Error())
			}
			return c, err
		},
	}
}

func newQueue(pool *redis.Pool, name string) (bool, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	newQ := Queue{Name: name, Offset: 0, Jobs: make([][]byte, 0)}
	// newQ.jobs = make([][]byte, 100)
	qBytes, _ := json.Marshal(newQ)
	n1, err := c.Do("SET", name, qBytes)

	if err != nil {
		return false, err
	}
	fmt.Println("Queue created with result: ", n1)
	return true, nil

}

func getQueue(pool *redis.Pool, name string) (*Queue, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	q := Queue{}
	qBytes, err := redis.Bytes(c.Do("GET", name))
	err = json.Unmarshal(qBytes, &q)

	if err != nil {
		return nil, err
	}
	fmt.Println("Retrieved queue: ", string(qBytes))
	return &q, nil
}

func createJob(pool *redis.Pool, name string, job []byte) (bool, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	// Get queue and append new job
	q, err := getQueue(pool, name)
	if err != nil {
		fmt.Printf("Could not retrieve queue")
		return false, err
	}
	q.Jobs = append(q.Jobs, job)

	// Marshal updated queue and store it
	qBytes, _ := json.Marshal(q)
	n1, err := c.Do("SET", name, qBytes)

	if err != nil {
		return false, err
	}
	fmt.Println("Queue updated with result: ", n1)
	return true, nil
}

func runJob(pool *redis.Pool, name string, handler func([]byte) []byte) ([]byte, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	// Get queue, extract job and update queue
	q, _ := getQueue(pool, name)
	var job []byte
	job, q.Jobs = q.Jobs[0], q.Jobs[1:]

	// Run the job using the handler function
	fmt.Println("Running job:", string(job))
	res := handler(job)

	// Marshal updated queue and store it
	qBytes, _ := json.Marshal(q)
	n1, err := c.Do("SET", name, qBytes)

	if err != nil {
		return nil, err
	}
	fmt.Println("Queue updated with result: ", n1)
	return res, nil
}

// newQ.Jobs = append(newQ.Jobs, []byte("testJob"))
// newQ.Jobs = append(newQ.Jobs, []byte("testJob"))
// x, a = a[0], a[1:]

func testRedis(pool *redis.Pool, command string, key string, value []byte) (bool, error) {
	// Init REDIS connection
	c := pool.Get()
	// Close connection once function is finished
	defer c.Close()

	_, err := c.Do(command, key, value)
	if err != nil {
		return false, err
	}

	return true, nil
}

// func main() {
// 	var pool = newPool()
// 	fmt.Println(newQueue(pool, "testQueue"))
// 	fmt.Println(getQueue(pool, "testQueue"))
// 	fmt.Println(createJob(pool, "testQueue", []byte("value")))
// 	fmt.Println(createJob(pool, "testQueue", []byte("value")))
// 	fmt.Println(getQueue(pool, "testQueue"))
// 	fmt.Println(runJob(pool, "testQueue"))
// 	fmt.Println(getQueue(pool, "testQueue"))
// 	// n1, _ := c.Do("SET", "key", []byte("value"))
// 	c := pool.Get()
// 	defer c.Close()
// 	// n1, _ := c.Do("SET", "key", []byte(""))
// 	n2, _ := redis.Bytes(c.Do("GET", "testQueue"))
// 	fmt.Println(string(n2))
// }
