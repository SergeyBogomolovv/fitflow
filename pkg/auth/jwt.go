package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignJWT(subject string, secret []byte, ttl time.Duration) (string, error) {
	iat := time.Now()
	exp := iat.Add(ttl)
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(iat),
		ExpiresAt: jwt.NewNumericDate(exp),
	}).SignedString(secret)
}

func VerifyJWT(tokenString string, secret []byte) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}
	return claims.Subject, nil
}

type AdminLoginKey struct{}
