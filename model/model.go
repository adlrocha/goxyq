package model

type Queue struct {
	offset int      `json:"offset"`
	name   string   `json:"name"`
	jobs   [][]byte `json:"jobs"`
}
