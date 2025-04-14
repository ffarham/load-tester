package models

import "time"

// Response is the reponse of the HTTP request.
type Response struct {
	StatusCode int
	Duration   time.Duration
}
