package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

func RateLimiter(maxRequests float64, burst int, ttl time.Duration) func(http.Handler) http.Handler {
	lml := tollbooth.NewLimiter(maxRequests, &limiter.ExpirableOptions{
		DefaultExpirationTTL: ttl,
	})

	lml.SetBurst(burst)
	lml.SetIPLookups([]string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"})
	lml.SetMessage(`{"error":"Too many requests. Please try again later."}`)
	lml.SetStatusCode(http.StatusTooManyRequests)
	lml.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rate limit exceeded: IP=%s", r.RemoteAddr)
		w.Header().Set("Retry-After", "60")
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpError := tollbooth.LimitByRequest(lml, w, r)
			if httpError != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(lml.GetMessage())) // Используем сообщение из лимитера
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
