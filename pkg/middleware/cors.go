package middleware

import "net/http"

// CORSMiddleware allows specific origin, methods, headers, and credentials
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Set allowed origin
		w.Header().Set("Access-Control-Allow-Origin", origin)
		// Allow credentials
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Allow the specified methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow the specified headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
