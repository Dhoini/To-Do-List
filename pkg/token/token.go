package token

import (
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
)

type JwtDate struct {
	UserId string
	Email  string
}

type JWTSecret struct {
	Secret string
}

func NewJWT(secret string) *JWTSecret {
	return &JWTSecret{
		Secret: secret,
	}
}

func (j *JWTSecret) GenerateToken(date JwtDate) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": date.UserId,
		"email":  date.Email,
	})

	secret, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		slog.Error(err.Error(), "can not sign token", err)
		return "", err
	}
	return secret, nil
}

func (j *JWTSecret) ParseToken(token string) (bool, *JwtDate) {
	slog.Info("Parsing Token", "token", token)
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		slog.Error(err.Error(), "can not parse token", err)
		return false, nil
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("invalid token claims")
		return false, nil
	}
	slog.Info("Token claims", "claims", claims)
	userID, ok := claims["userId"].(string)
	if !ok {
		slog.Error("invalid userId in token claims", "actual_value", claims["userID"])
		return false, nil
	}
	email, ok := claims["email"].(string)
	if !ok && claims["email"] != nil {
		slog.Error("invalid email in token claims", "actual_value", claims["email"])
		return false, nil
	}
	return t.Valid, &JwtDate{
		UserId: userID,
		Email:  email,
	}
}
