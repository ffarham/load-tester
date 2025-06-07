package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Url           string
	Method        string
	Payload       interface{}
	WorkerId      int
	IsEventSource bool
}

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
		if v.Method == http.MethodGet && v.IsEventSource {
			setEventSourceHeaders(req)
		}

		start := time.Now()
		res, err := client.Do(req)
		duration := time.Since(start)
		if err != nil {
			// when connection is refused, returning 503 Service Unavailable response
			output <- Response{
				StatusCode: 503,
				Duration:   duration,
			}
			continue
		}

		defer res.Body.Close()

		if v.Method == http.MethodGet && v.IsEventSource {
			handleEventSourceResponse(res, output)
			continue
		}

		output <- Response{
			StatusCode: res.StatusCode,
			Duration:   duration,
		}
	}
}

func handleEventSourceResponse(res *http.Response, output chan<- Response) {
	sseStreamStartTime := time.Now()
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		scanner.Text()
	}
	output <- Response{
		StatusCode: res.StatusCode,
		Duration:   time.Since(sseStreamStartTime),
	}
}

func setEventSourceHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
}
