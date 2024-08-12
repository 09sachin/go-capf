package routes

import (
	"net/http"
	"sync"
	"time"
    "io"
    "bytes"
    "encoding/json"
	"github.com/09sachin/go-capf/controllers"
	"github.com/gorilla/mux"
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
    if len(rl.requests) > 1000 {
        rl.requests = make(map[string]map[time.Time]int) // Clear the entire requests map
    }
}


type RequestBody struct {
	ForceID string `json:"force_id"`
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
            ip := r.RemoteAddr 

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
                controllers.ErrorLogger.Println("Rate limit exceeded - Attack")
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


// RateLimitMiddleware creates a middleware for rate limiting requests based on  force Id
func (rl *RateLimiter) LoginRateLimitMiddleware(next http.Handler, limit int, duration time.Duration, cleanupInterval time.Duration) mux.MiddlewareFunc {
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
            body, err := io.ReadAll(r.Body)
            if err != nil {
                http.Error(w, "Error reading request body", http.StatusBadRequest)
                return
            }
            r.Body = io.NopCloser(bytes.NewReader(body))

            // Unmarshal the request body into the struct
            var requestBody RequestBody
            err = json.Unmarshal(body, &requestBody)
            if err != nil {
                http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
                return
            }

            // Access the force_id field
            forceID := requestBody.ForceID
            ip := forceID 

            rl.mu.Lock()
            defer rl.mu.Unlock()

            if _, ok := rl.requests[ip]; !ok {
                rl.requests[ip] = make(map[time.Time]int)
            }

            for t := range rl.requests[ip] {
                if time.Since(t) > duration {
                    delete(rl.requests[ip], t)
                }
            }

            // Check if the request count exceeds the limit
            if len(rl.requests[ip]) >= limit {
                controllers.ErrorLogger.Printf("Rate limit exceeded : %s", ip)
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
    loginLimiter := NewRateLimiter()
    restrictedRoute := route.PathPrefix("/login").Subrouter()
    dataRoute := route.PathPrefix("/data").Subrouter()
    dataRoute.Use(loginLimiter.RateLimitMiddleware(route, 100, time.Minute,   10 * time.Minute))
    restrictedRoute.Use(rateLimiter.LoginRateLimitMiddleware(route, 10, 5*time.Minute,   10 * time.Minute))
	restrictedRoute.HandleFunc("/send-otp", controllers.SendOtp).Methods("POST")
	restrictedRoute.HandleFunc("/otp-login", controllers.OtpLogin).Methods("POST")
	dataRoute.HandleFunc("/dashboard-data", controllers.DashboardData).Methods("GET")
	dataRoute.HandleFunc("/user-details", controllers.UserDetails).Methods("GET")
	dataRoute.HandleFunc("/hospital-list", controllers.Hospitals).Methods("GET")
	dataRoute.HandleFunc("/filter-hospital", controllers.FilterHospital).Methods("GET")
	dataRoute.HandleFunc("/queries", controllers.Queries).Methods("GET")
	dataRoute.HandleFunc("/track-case", controllers.TrackCases).Methods("GET")
	dataRoute.HandleFunc("/claims", controllers.UserClaims).Methods("GET")
    dataRoute.HandleFunc("/claims/get-fields", controllers.GetUpdateClaimsFieldsAPI).Methods("GET")
    dataRoute.HandleFunc("/claims/update-api", controllers.UpdateClaimsAPI).Methods("POST")

	return route
}
