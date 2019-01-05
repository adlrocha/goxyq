package queue

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SimpleJob serialization
type SimpleJob struct {
	Command string `json:"command"`
	Arg1    int    `json:"arg1"`
	Arg2    int    `json:"arg2"`
}

// simpleJobHandler for a testJob
// Handler receive []byte and returns []byte. Is the developer responsabilityç
// to handle data correctly between jobs and functions
func simpleJobHandler(job []byte) ([]byte, error) {
	j := SimpleJob{}
	json.Unmarshal(job, &j)
	res := j.Arg1 + j.Arg2
	fmt.Println("Running job with command: ", j.Command)
	return []byte(strconv.Itoa(res)), nil
}
