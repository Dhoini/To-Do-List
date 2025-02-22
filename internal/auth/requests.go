package auth

import (
	"ToDo/internal/user"
	"ToDo/pkg/req"
	"ToDo/pkg/res"
	token2 "ToDo/pkg/token"
	"errors"
	"net/http"
)

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		userId, err := h.AuthService.Register(r.Context(), body.Email, body.Password, body.Name)
		if err != nil {
			if errors.Is(err, user.ErrUserAlreadyExists) {
				res.JsonResponse(w, res.ErrorResponse{Error: err.Error()}, http.StatusConflict) //
			} else {
				res.JsonResponse(w, res.ErrorResponse{Error: "internal server error"}, http.StatusInternalServerError) //
			}
			return
		}
		token, err := token2.NewJWT(h.Config.Auth.Secret).GenerateToken(token2.JwtDate{
			UserId: userId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			Token: token,
		}
		res.JsonResponse(w, data, http.StatusOK)

	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}

		email, err := h.AuthService.Login(r.Context(), body.Email, body.Password)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				res.JsonResponse(w, res.ErrorResponse{Error: "invalid credentials"}, http.StatusUnauthorized)
			} else {
				res.JsonResponse(w, res.ErrorResponse{Error: "internal server errror"}, http.StatusInternalServerError)
			}
			return
		}

		token, err := token2.NewJWT(h.Config.Auth.Secret).GenerateToken(token2.JwtDate{
			Email: email.Email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := LoginResponse{
			Token: token,
		}
		res.JsonResponse(w, data, http.StatusOK)

	}
}
