package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		nextHandler.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		statusCode := wrapper.StatusCode

		if statusCode >= http.StatusBadRequest {
			log.Printf("[ERROR] %d %s %s %v", statusCode, r.Method, r.URL.Path, duration) // Логируем как ошибку
		} else {
			log.Printf("[INFO] %d %s %s %v", statusCode, r.Method, r.URL.Path, duration) // Логируем как информацию
		}
	})
}
