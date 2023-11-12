package ratelimiter

import (
	"log"
	"net"
	"net/http"
)

var rateLimiter = NewIPRateLimiter(1, 5)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("🚀 Server is running on port 8888")
	log.Fatal(http.ListenAndServe(":8888", RateLimit(mux)))
}

// Middleware ...
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		if remoteAddr == "" {
			remoteAddr = r.Header.Get("X-Forwarded-For")
		}
		if remoteAddr == "" {
			remoteAddr = r.Header.Get("X-Real-IP")
		}
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if limiter := rateLimiter.LimiterByIP(ip); !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
