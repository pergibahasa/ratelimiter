package ratelimiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPRateLimiter ...
type IPRateLimiter struct {
	// Registered as a visitor by IP Address
	visitors map[string]*visitor

	mu *sync.RWMutex

	// r defines the maximum frequency of some events. r is represented as number of events per second. A zero Limit allows no events.
	//
	// permits you to consume an average of r tokens per second, with a maximum of b tokens in any single 'burst'.
	r rate.Limit

	// maximum burst size
	b int
}

var (
	ipRateLimiterConfig *IPRateLimiter
	ipRateLimiterOnce   sync.Once
)

// NewIPRateLimiter ...
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	ipRateLimiterOnce.Do(func() {
		ipRateLimiterConfig = &IPRateLimiter{
			visitors: make(map[string]*visitor),
			mu:       &sync.RWMutex{},
			r:        r,
			b:        b,
		}

		go ipRateLimiterConfig.cleanupVisitors()
	})
	return ipRateLimiterConfig
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.visitors[ip] = &visitor{limiter: limiter, lastSeenAt: time.Now()}

	return limiter
}

// LimiterByIP returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) LimiterByIP(ip string) *rate.Limiter {
	i.mu.Lock()

	var (
		v  *visitor
		ok bool
	)
	if v, ok = i.visitors[ip]; !ok {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	v.lastSeenAt = time.Now()

	i.mu.Unlock()

	return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func (r *IPRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		r.mu.Lock()
		for ip, v := range r.visitors {
			if time.Since(v.lastSeenAt) > 3*time.Minute {
				delete(r.visitors, ip)
			}
		}
		r.mu.Unlock()
	}
}
