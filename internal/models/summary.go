package models

import (
	"fmt"
	"math"
)

type Summary struct {
	TotalRequests               int
	SuccessfulRequests          int
	ThrottledRequests           int
	FailedRequests              int
	ConnectionFailures          int
	TotalTime                   float64
	RequestsPerSecond           float64
	SuccessfulRequestsPerSecond float64
	MinResponseTime             float64
	AverageResponseTime         float64
	MaxResponseTime             float64
}

func (v *Summary) String() string {
	return fmt.Sprintf(`
	
	---Load Test Results---
Total Requests:        	%d
Successful requests:    %d
Failed Requests:        %d
Throttled Requests:     %d
Refused Connections:    %d

[Successful Requests]
Total Time:             %s%s
Min Response Time:      %s%s
Average Response Time:  %s%s
Max Response Time:      %s%s
Requests Per Second:    %s
		
		`,
		v.TotalRequests, v.SuccessfulRequests, v.FailedRequests, v.ThrottledRequests,
		v.ConnectionFailures,
		getTimeValue(v.TotalTime), getTimeUnit(v.TotalTime),
		getTimeValue(v.MinResponseTime), getTimeUnit(v.MinResponseTime),
		getTimeValue(v.AverageResponseTime), getTimeUnit(v.AverageResponseTime),
		getTimeValue(v.MaxResponseTime), getTimeUnit(v.MaxResponseTime),
		getTimeValue(v.RequestsPerSecond),
	)
}

func getTimeValue(value float64) string {
	if math.IsNaN(value) {
		return "-"
	}
	return fmt.Sprintf("%f", value)
}

func getTimeUnit(value float64) string {
	if math.IsNaN(value) {
		return ""
	}
	return "s"
}
