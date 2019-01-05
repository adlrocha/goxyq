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
	res, err := newQueue(pool, "a1")
	if !res || err != nil {
		t.FailNow()
	}

	a, _ := getQueue(pool, "a1")
	if a.Name != "a1" {
		t.FailNow()
	}
}

func TestJob(t *testing.T) {
	var pool = newPool()
	// Create queue
	res, err := newQueue(pool, "a1")
	// Create job
	fmt.Println("[*] Creating new job")
	job := SimpleJob{"sum", 50, 1000000}
	jobBytes, _ := json.Marshal(job)
	res, err = createJob(pool, "a1", jobBytes)
	if !res || err != nil {
		t.FailNow()
	}

	fmt.Println("[*] Running job")
	jobRes, _ := runJob(pool, "a1", simpleJobHandler)
	fmt.Println("Result from execution", string(jobRes))
	if !res || err != nil {
		t.FailNow()
	}
}

// func TestTrust(t *testing.T) {
// 	// 	scc := new(SmartContract)
// 	// 	stub := shim.NewMockStub("trust", scc)

// 	// 	checkInit(t, stub, nil)

// 	// 	asset := []byte(`{"id":"as1","owner":"admin","data":{"metric1":"a","metric2":"b"}}`)
// 	// 	metadata := []byte(`{"metric1":"c","metric2":"b", "metric3":"c"}`)
// 	// 	checkInvoke(t, stub, [][]byte{[]byte("create"), asset})
// 	// 	checkGetAsset(t, stub, "as1")
// 	// 	checkInvoke(t, stub, [][]byte{[]byte("update"), []byte("as1"), metadata})
// 	// 	checkGetAsset(t, stub, "as1")
// 	// 	// checkHistory(t, stub, "as1")
// }
