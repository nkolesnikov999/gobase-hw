package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{secret: secret}
}

func (j *JWT) Create(subject string) (string, error) {
	claims := jwt.MapClaims{
		"sub": subject,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return s, nil
}
