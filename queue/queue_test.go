package queue

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestQueue Tests to create a queue and get it.
func TestQueue(t *testing.T) {
	var pool = NewPool()
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

func TestSingleJob(t *testing.T) {
	var pool = NewPool()
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
	jobRes, err := WaitAndRunJob(pool, "a1", jobBytes)
	if err != nil {
		fmt.Println("Error in job run: ", err)
		t.FailNow()
	}
	fmt.Println("Result from job jun run", jobRes)

	fmt.Println("[*] Verify status of the queue")
	a, _ := GetQueue(pool, "a1")
	if a.Name != "a1" {
		t.FailNow()
	}
}

func TestMultipleJobs(t *testing.T) {
	var pool = NewPool()
	// Create queue
	res, err := NewQueue(pool, "a1")
	// Create job
	fmt.Println("[*] Creating two new jobs")
	job1 := SimpleJob{"sum", 50, 1000000}
	job2 := SimpleJob{"mul", 500, 11000000}
	jobBytes1, _ := json.Marshal(job1)
	jobBytes2, _ := json.Marshal(job2)
	res, err = CreateJob(pool, "a1", jobBytes1)
	if !res || err != nil {
		t.FailNow()
	}
	res, err = CreateJob(pool, "a1", jobBytes2)
	if !res || err != nil {
		t.FailNow()
	}

	fmt.Println("[*] Running job 1")
	jobRes, err := WaitAndRunJob(pool, "a1", jobBytes1)
	if err != nil {
		fmt.Println("Error in job run: ", err)
		t.FailNow()
	}
	fmt.Println("Result from job run", jobRes)

	fmt.Println("[*] Running job 2")
	jobRes, err = WaitAndRunJob(pool, "a1", jobBytes2)
	if err != nil {
		fmt.Println("Error in job run: ", err)
		t.FailNow()
	}
	fmt.Println("Result from job run", jobRes)

	fmt.Println("[*] Verify status of the queue")
	a, _ := GetQueue(pool, "a1")
	if a.Name != "a1" {
		t.FailNow()
	}
}
