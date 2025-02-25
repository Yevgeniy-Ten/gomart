package session

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

func CreateToken(userID int) (string, error) {
	claims := &Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

const AuthValidLength = 7

// GetUserID Bearer ${token}
func GetUserID(authString string) (int, error) {
	if len(authString) < AuthValidLength {
		return 0, errors.New("invalid token")
	}
	tokenStr := authString[AuthValidLength:]
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("session is not valid")
	}

	return claims.UserID, nil
}
