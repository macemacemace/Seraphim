package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const accessTokenTTL = 24 * time.Hour

func newAccessToken(userID uuid.UUID, secret []byte) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func parseAccessToken(tokenString string, secret []byte) (uuid.UUID,
	error) {
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(claims.Subject)

}
