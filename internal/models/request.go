package models

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Request is the http request to be executed by a worker thread.
type Request struct {
	Url      string
	Method   string
	Payload  interface{}
	WorkerId int
}

// Execute executes the http request and sends the response to the output channel.
func (v *Request) Execute(client *http.Client, n int, output chan<- Response) {
	for index := range n {
		body, err := json.Marshal(v.Payload)
		if err != nil {
			log.Fatalf(
				"Error marshalling payload by worker %d for request %d\n%v",
				v.WorkerId, index, err.Error(),
			)
		}
		bodyReader := bytes.NewReader(body)
		req, err := http.NewRequest(v.Method, v.Url, bodyReader)
		if err != nil {
			log.Fatalf(
				"Error creating request %d for worker %d\n%v",
				index, v.WorkerId, err.Error(),
			)
		}

		start := time.Now()
		res, err := client.Do(req)
		duration := time.Since(start)
		if err != nil {
			log.Fatalf(
				"Error executing request %d by worker %d\n%v",
				index, v.WorkerId, err.Error(),
			)
		}

		res.Body.Close()

		output <- Response{
			StatusCode: res.StatusCode,
			Duration:   duration,
		}
	}
}
