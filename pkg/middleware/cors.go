package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			headers.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			headers.Set("Access-Control-Max-Age", "86400")
			headers.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, content-length")
			return
		}
		next.ServeHTTP(w, r)
	})
}
