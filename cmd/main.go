package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/ffarham/load-tester/internal/models"
	"github.com/ffarham/load-tester/internal/utils"
)

const (
	defaultWorkers  = 1
	defaultRequests = 1
	defaultTimeout  = 3
	defaultRequired = "REQUIRED"
)

// example usage:
// go run cmd/main.go -m=GET -url=http://localhost:8080/ -w=10 -r=100 -t=10
// `m`: required flag for specifying the HTTP method (GET, POST, e.t.c.).
// `url`: required flag for specifying the URL to load test.
// `w`: flag for specifying the number of concurrent workers. Default is 1.
// `r`: flag for specifying the total number of requests to send. Default is 1.
// `t`: flag for specifying the HTTP request timeout in seconds. Default is 3.
// `f`: flag for specifying the filename to load the json payload from.
// `i`: flag for specifying whether to make an insecure request. Default is false.s
// `sse`: flag to specify whether the request is a SSE request.
// Payload is `nil` for post reqests when this flag is not provided, as well as
// for other requests. Default is "".
func main() {

	method := flag.String("m", defaultRequired, "HTTP method to use")
	url := flag.String("url", defaultRequired, "URL to load test")
	workers := flag.Int("w", defaultWorkers, "Number of concurrent workers")
	requests := flag.Int("r", defaultRequests, "Total number of requests to send")
	insecure := flag.Bool("i", false, "Skip TLS verification")
	sse := flag.Bool("sse", false, "Indicates if the request is for Server-Sent Events (text/event-stream)")
	timeout := flag.Int("t", defaultTimeout, "Timeout in seconds")
	filename := flag.String("f", "", "Filename to load payload from")
	flag.Parse()

	err := validateCmdLineInput(*method, *url, *workers, *requests, *timeout)
	if err != nil {
		log.Fatalf("Error validating command line inputs.\n%v", err)
	}

	payload := getPayload(*method, *filename)

	log.Printf(
		"Starting load tester with %d workers makings %d %s requests to %s",
		*workers, *requests, *method, *url,
	)

	// http.Transport object that skips TLS certificate verification if insecure flag is set
	var transport http.Transport
	if *insecure {
		transport = http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	var client *http.Client
	if *sse {
		client = &http.Client{
			Transport: &transport,
		}
	} else {
		client = &http.Client{
			Transport: &transport,
			Timeout:   time.Duration(*timeout) * time.Second,
		}
	}

	// create a buffered channel to hold the responses
	responses := make(chan models.Response, *requests)
	var wg sync.WaitGroup

	reqsPerWorker := *requests / *workers
	remainder := *requests % *workers

	for i := 0; i < *workers; i++ {
		wg.Add(1)

		// handle the remainder
		workerRequests := reqsPerWorker
		if i < remainder {
			workerRequests++
		}

		go func() {
			defer wg.Done()
			request := &models.Request{
				Url:           *url,
				Method:        *method,
				Payload:       payload,
				WorkerId:      i,
				IsEventSource: *sse,
			}
			request.Execute(client, workerRequests, responses)
		}()
	}

	// wait for all workers to finish before closing the channel
	go func() {
		wg.Wait()
		close(responses)
	}()

	summary := summarise(responses)
	log.Printf("%s", summary)
}

func validateCmdLineInput(method, url string, workers, requests, timeout int) error {
	if len(method) == 0 {
		return errors.New("HTTP method is required")
	}
	if len(url) == 0 {
		return errors.New("URL is required")
	}
	if workers <= 0 {
		return errors.New("number of workers must be at least 1")
	}
	if requests <= 0 {
		return errors.New("number of requests must be at least 1")
	}
	if timeout <= 0 {
		return errors.New("timeout must be at least 1 second")
	}
	return nil
}

func getPayload(method, filename string) (payload interface{}) {
	if method == "POST" && len(filename) > 0 {
		jsonData, err := utils.ReadJsonFile(filename)
		if err != nil {
			log.Fatalf("Error loading payload from file %s\n%v", filename, err)
		}
		payload = jsonData
		log.Printf("Loaded payload from file %s", filename)
	}

	return
}

func summarise(responses chan models.Response) *models.Summary {

	summary := &models.Summary{
		MinResponseTime: math.MaxFloat64,
	}

	for response := range responses {
		summary.TotalRequests++

		if response.StatusCode == http.StatusServiceUnavailable {
			summary.ConnectionFailures++
			continue
		}

		if response.StatusCode == http.StatusOK {
			summary.SuccessfulRequests++
			summary.TotalTime += response.Duration.Seconds()
			summary.MinResponseTime = min(summary.MinResponseTime, response.Duration.Seconds())
			summary.MaxResponseTime = max(summary.MaxResponseTime, response.Duration.Seconds())
		} else if response.StatusCode == http.StatusTooManyRequests {
			summary.ThrottledRequests++
		} else {
			summary.FailedRequests++
		}
	}
	summary.RequestsPerSecond = utils.SafeDiv(float64(summary.SuccessfulRequests), summary.TotalTime)
	summary.AverageResponseTime = utils.SafeDiv(summary.TotalTime, float64(summary.SuccessfulRequests))

	// edge case: no connection was established
	if summary.MinResponseTime == math.MaxFloat64 {
		summary.TotalTime = math.NaN()
		summary.MinResponseTime = math.NaN()
		summary.MaxResponseTime = math.NaN()
	}

	return summary
}
