package middleware

import (
	"ToDo/configs"
	token2 "ToDo/pkg/token"
	"context"
	"net/http"
	"strings"
)

type key string

const ContextUserIDKey key = "userID"

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	if err != nil {
		panic(err)
	}
}

func IsAuthenticated(config *configs.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeUnauthorized(w)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == " " {
				writeUnauthorized(w)
				return
			}
			isValid, data := token2.NewJWT(config.Auth.Secret).ParseToken(token)
			if !isValid {
				writeUnauthorized(w)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserIDKey, data.UserId)
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}
}
