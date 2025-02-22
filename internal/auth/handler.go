package auth

import (
	"ToDo/configs"
	"ToDo/pkg/di"
	"ToDo/pkg/middleware"
	"net/http"
)

// AuthHandler — структура хендлера, тоже используем интерфейс
type AuthHandler struct {
	Config      *configs.Config
	AuthService di.IAuthService // Заменяем *AuthService на интерфейс
}

type AuthHandlerDeps struct {
	Config      *configs.Config
	AuthService di.IAuthService
}

func NewAuthHandler(router *http.ServeMux, deps *AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	middlewares := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
		middleware.RateLimiter(deps.Config.RateLimit.MaxRequests, deps.Config.RateLimit.Burst, deps.Config.RateLimit.TTL),
	)

	router.Handle("POST /auth/login", middlewares(handler.Login()))
	router.Handle("POST /auth/register", middlewares(handler.Register()))
}
