package models

import "fmt"

// Summary is the summary of the load test results.
type Summary struct {
	TotalRequests               int
	SuccessfulRequests          int
	ThrottledRequests           int
	OtherFailedRequests         int
	TotalTime                   float64
	RequestPerSecond            float64
	SuccessfulRequestsPerSecond float64
	MinResponseTime             float64
	AverageResponseTime         float64
	MaxResponseTime             float64
}

// String implements the fmt.Stringer interface for the Summary struct.
func (v *Summary) String() string {
	return fmt.Sprintf(`
	
	---Load Test Results---
Total Requests:        	%d
Successful requests:    %d
Throttled Requests:     %d
Other Failed Requests:  %d
Requests Per Second:     %f

Total Time:             %fs
Min Response Time:      %fs
Average Response Time:  %fs
Max Response Time:      %fs
		
		`,
		v.TotalRequests, v.SuccessfulRequests, v.ThrottledRequests, v.OtherFailedRequests,
		v.RequestPerSecond, v.TotalTime, v.MinResponseTime, v.AverageResponseTime,
		v.MaxResponseTime,
	)
}
