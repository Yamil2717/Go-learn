package main

import (
	"encoding/json"
	"time"

	"github.com/Yamil2717/todo"
)

type todoResponse struct {
	Results todo.List
}

func (tr todoResponse) MarshalJSON() ([]byte, error) {
	response := struct {
		Results      todo.List `json:"results"`
		Date         time.Time `json:"date"`
		TotalResults int       `json:"total_results"`
	}{
		Results:      tr.Results,
		Date:         time.Now(),
		TotalResults: len(tr.Results),
	}

	return json.Marshal(response)
}
