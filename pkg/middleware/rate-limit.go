package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

/* Improvement Points
Check the X-Forwarded-For or X-Real-IP headers for the IP address, if you are running your server behind a reverse proxy.
Port the code to a standalone package.
Make the rate limiter and cleanup settings configurable at runtime.
Remove the reliance on global variables, so that different rate limiters can be created with different settings.
Switch to a sync.RWMutex to help reduce contention on the map.
*/

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimit struct {
	visitors map[string]*visitor
	limit rate.Limit
	burst int
}


// Change the the map to hold values of the type visitor.
var mu sync.Mutex

func NewRateLimitMiddleware(rateLimit rate.Limit, rateBurst int) *RateLimit{
	return &RateLimit{make(map[string]*visitor), rateLimit, rateBurst}
}

// Run a background goroutine to remove old entries from the visitors map.
func (rateLimit *RateLimit) init() {
	go rateLimit.cleanupVisitors()
}

func (rateLimit *RateLimit) getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := rateLimit.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rateLimit.limit, rateLimit.burst)
		// Include the current time when creating a new visitor.
		rateLimit.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func (rateLimit *RateLimit) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		defer mu.Unlock()
		for ip, v := range rateLimit.visitors {
			if time.Now().Sub(v.lastSeen) > 3*time.Minute {
				delete(rateLimit.visitors, ip)
			}
		}
	}
}

func (rateLimit *RateLimit) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := rateLimit.getVisitor(ip)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}