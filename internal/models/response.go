package models

import "time"

type Response struct {
	StatusCode int
	Duration   time.Duration
}
