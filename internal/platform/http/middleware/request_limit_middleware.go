package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"finsolvz-backend/internal/utils"
)

// RequestLimitMiddleware adds request size limits and timeouts
func RequestLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Limit request body size to 10MB
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		
		// Set request timeout context
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		
		// Add timeout header
		w.Header().Set("X-Request-Timeout", "30s")
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RateLimitMiddleware implements basic rate limiting (in-memory)
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	type client struct {
		requests int
		window   time.Time
	}
	
	clients := make(map[string]*client)
	var mutex sync.RWMutex
	
	// Cleanup old entries every minute
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				mutex.Lock()
				now := time.Now()
				for ip, c := range clients {
					if now.Sub(c.window) > time.Minute {
						delete(clients, ip)
					}
				}
				mutex.Unlock()
			}
		}
	}()
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}
			
			mutex.Lock()
			c, exists := clients[ip]
			if !exists {
				c = &client{
					requests: 0,
					window:   time.Now(),
				}
				clients[ip] = c
			}
			
			// Reset window if it's been more than a minute
			if time.Since(c.window) > time.Minute {
				c.requests = 0
				c.window = time.Now()
			}
			
			c.requests++
			currentRequests := c.requests
			mutex.Unlock()
			
			if currentRequests > requestsPerMinute {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "60")
				
				utils.RespondJSON(w, http.StatusTooManyRequests, map[string]string{
					"error": "Rate limit exceeded",
					"message": "Too many requests, please try again later",
				})
				return
			}
			
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", requestsPerMinute-currentRequests))
			
			next.ServeHTTP(w, r)
		})
	}
}