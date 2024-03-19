package routes

import (
	"github.com/09sachin/go-capf/controllers"
    "github.com/09sachin/go-capf/config"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)


type RateLimiter struct {
    mu       sync.Mutex
    requests map[string]map[time.Time]int // Map of IP addresses to request counts and timestamps
}


func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        requests: make(map[string]map[time.Time]int),
    }
}

// // Print the requests map
// func (rl *RateLimiter) PrintRequestsMap() {
// 	for ip, timestamps := range rl.requests {
// 		log.Printf("IP: %s\n", ip)
// 		for t, count := range timestamps {
// 			log.Printf("Timestamp: %s, Count: %d\n", t.Local(), count)
// 		}
// 	}
// }

func (rl *RateLimiter) cleanupOldEntries() {
    rl.mu.Lock()
    defer rl.mu.Unlock()
	// rl.PrintRequestsMap()
    rl.requests = make(map[string]map[time.Time]int) // Clear the entire requests map
}


// RateLimitMiddleware creates a middleware for rate limiting requests based on IP address.
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler, limit int, duration time.Duration, cleanupInterval time.Duration) mux.MiddlewareFunc {
    go func() {
        ticker := time.NewTicker(cleanupInterval)
        defer ticker.Stop()
        for {
            <-ticker.C
            rl.cleanupOldEntries()
        }
    }()

	return func(handler http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr // Assuming IP is in the RemoteAddr field

            rl.mu.Lock()
            defer rl.mu.Unlock()

            // Initialize the map entry for this IP if it doesn't exist
            if _, ok := rl.requests[ip]; !ok {
                rl.requests[ip] = make(map[time.Time]int)
            }

            // Remove old entries from the map
            for t := range rl.requests[ip] {
                if time.Since(t) > duration {
                    delete(rl.requests[ip], t)
                }
            }

            // Check if the request count exceeds the limit
            if len(rl.requests[ip]) >= limit {
                controllers.ErrorLogger.Println("Rate limit exceeded")
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }

            // Increment the request count for the current time
            rl.requests[ip][time.Now()]++

            // Call the next handler
            handler.ServeHTTP(w, r)
        })
    }
}


func Init() *mux.Router {
	route := mux.NewRouter()
	rateLimiter := NewRateLimiter()
	route.HandleFunc("/send-otp", controllers.SendOtp).Methods("POST")
	route.HandleFunc("/otp-login", controllers.OtpLogin).Methods("POST")
	route.HandleFunc("/dashboard-data", controllers.DashboardData).Methods("GET")
	route.HandleFunc("/user-details", controllers.UserDetails).Methods("GET")
	route.HandleFunc("/hospital-list", controllers.Hospitals).Methods("GET")
	route.HandleFunc("/filter-hospital", controllers.FilterHospital).Methods("GET")
	route.HandleFunc("/queries", controllers.Queries).Methods("GET")
	route.HandleFunc("/track-case", controllers.TrackCases).Methods("GET")
	route.HandleFunc("/claims", controllers.UserClaims).Methods("GET")
    route.HandleFunc("/api/devops/nhvb3bhb4/deploy", config.Deployments).Methods("POST")
	route.Use(rateLimiter.RateLimitMiddleware(route, 10, time.Second,   10 * time.Second))  

	return route
}
