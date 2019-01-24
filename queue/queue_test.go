package queue

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestQueue Tests to create a queue and get it.
func TestQueue(t *testing.T) {
	var pool = newPool()
	fmt.Println("[*] Creating new queue")
	res, err := NewQueue(pool, "a1")
	if !res || err != nil {
		t.FailNow()
	}

	a, _ := GetQueue(pool, "a1")
	if a.Name != "a1" {
		t.FailNow()
	}
}

func TestJob(t *testing.T) {
	var pool = newPool()
	// Create queue
	res, err := NewQueue(pool, "a1")
	// Create job
	fmt.Println("[*] Creating new job")
	job := SimpleJob{"sum", 50, 1000000}
	jobBytes, _ := json.Marshal(job)
	res, err = CreateJob(pool, "a1", jobBytes)
	if !res || err != nil {
		t.FailNow()
	}

	fmt.Println("[*] Running job")
	jobRes, _ := RunJob(pool, "a1", simpleJobHandler)
	fmt.Println("Result from execution", string(jobRes))
	if !res || err != nil {
		t.FailNow()
	}
}
