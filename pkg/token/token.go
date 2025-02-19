package token

import (
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
)

type JwtDate struct {
	Email string
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
		"email": date.Email,
	})

	secret, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		slog.Error(err.Error(), "can not sign token")
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
		slog.Error(err.Error(), "can not parse token")
		return false, nil
	}

	email := t.Claims.(jwt.MapClaims)["email"].(string)
	return t.Valid, &JwtDate{
		Email: email,
	}
}
