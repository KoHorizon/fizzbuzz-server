package middleware

import (
	"net/http"
)

// CORSMiddleware adds CORS headers to allow cross-origin requests
// This is particularly useful for development and testing with tools like Swagger Editor
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow requests from any origin (for development/testing)
			// In production, you might want to restrict this to specific domains
			w.Header().Set("Access-Control-Allow-Origin", "*")

			// Allow common HTTP methods
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			// Allow common headers
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

			// Allow credentials (cookies, authorization headers)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Cache preflight requests for 1 hour
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}
