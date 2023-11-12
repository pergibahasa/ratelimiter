package ratelimiter

import (
	"time"

	"golang.org/x/time/rate"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter    *rate.Limiter
	lastSeenAt time.Time
}
