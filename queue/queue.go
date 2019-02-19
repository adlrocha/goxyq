package queue

import (
	"encoding/json"
	"fmt"

	"github.com/adlrocha/goxyq/log"
	"github.com/gomodule/redigo/redis"
)

// Queue represents a proxy queue for requests
type Queue struct {
	Name   string   `json:"name"`
	Offset int      `json:"offset"`
	Jobs   [][]byte `json:"jobs"`
}

type redisPool redis.Pool

//NewPool Creates a new REDIS Pool
func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				log.Fatalf("Could not connect to REDIS DB\n")
				panic(err.Error())
			}
			return c, err
		},
	}
}

// NewQueue creates a new job queue
func NewQueue(pool *redis.Pool, name string) (bool, error) {
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
	log.Infof("Queue created with result: %v", n1)
	return true, nil

}

// GetQueue gets the status of a queue
func GetQueue(pool *redis.Pool, name string) (*Queue, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	q := Queue{}
	qBytes, err := redis.Bytes(c.Do("GET", name))
	err = json.Unmarshal(qBytes, &q)

	if err != nil {
		return nil, err
	}
	log.Infof("Retrieved queue: %v\n", string(qBytes))
	return &q, nil
}

// CreateJob creates a new job for the queue
func CreateJob(pool *redis.Pool, name string, job []byte) (bool, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	// Get queue and append new job
	q, err := GetQueue(pool, name)
	if err != nil {
		log.Warnf("Could not retrieve queue: %v \n", name)
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

// EmptyQueue empty the queue of jobs
func EmptyQueue(pool *redis.Pool, name string) (bool, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	// Get queue and append new job
	q, err := GetQueue(pool, name)
	if err != nil {
		fmt.Printf("Could not retrieve queue: %v \n", name)
		return false, err
	}
	q.Jobs = make([][]byte, 0)

	// Marshal updated queue and store it
	qBytes, _ := json.Marshal(q)
	n1, err := c.Do("SET", name, qBytes)

	if err != nil {
		return false, err
	}
	log.Warnf("Could not retrieve queue: %v \n", n1)
	return true, nil
}

// Auxiliary function to check equality between byte[]
func bytesEqual(a []byte, b []byte) (res bool) {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// RunJobWithHandler runs the next job of a queue
func RunJobWithHandler(pool *redis.Pool, name string, handler func([]byte) ([]byte, error)) ([]byte, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	// Get queue, extract job and update queue
	q, _ := GetQueue(pool, name)
	var job []byte
	job, q.Jobs = q.Jobs[0], q.Jobs[1:]

	// Run the job using the handler function
	log.Infof("Running job: %v\n", string(job))
	res, err := handler(job)
	if err != nil {
		log.Errorf("Error executing job handler")
		return nil, err
	}

	// Marshal updated queue and store it
	qBytes, _ := json.Marshal(q)
	n1, err := c.Do("SET", name, qBytes)

	if err != nil {
		return nil, err
	}
	log.Infof("Queue updated with result: %v\n", n1)
	return res, nil
}

// WaitAndRunJob runs the next job of a queue
func WaitAndRunJob(pool *redis.Pool, name string, job []byte) (bool, error) {
	// Create and close connection when finished
	c := pool.Get()
	defer c.Close()

	q, err := GetQueue(pool, name)
	if q == nil || err != nil {
		return false, err
	}

	// While not job's turn, wait
	isTurn := bytesEqual(job, q.Jobs[0])
	for !isTurn {
		// Get queue, and verify if is my turn
		q, _ = GetQueue(pool, name)
		isTurn = bytesEqual(job, q.Jobs[0])
	}

	// Update Job list
	job, q.Jobs = q.Jobs[0], q.Jobs[1:]

	// // Run the job using the handler function
	// fmt.Println("Running job:", string(job))

	// Marshal updated queue and store it
	qBytes, _ := json.Marshal(q)
	n1, err := c.Do("SET", name, qBytes)

	if err != nil {
		return false, err
	}
	log.Infof("Queue updated with result: %v\n", n1)
	return true, nil
}

// Dummy function to test REDIS connection
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
