package jwt

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const CookieName = "token"

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

func GetAuthCookie(req *http.Request) (string, error) {
	cookie, err := req.Cookie(CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: uuid.New().String(),
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func fromString(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func GetUserID(tokenString string) (string, error) {
	claims, err := fromString(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

func IsTokenExpired(tokenString string) bool {
	claims, err := fromString(tokenString)
	if err != nil {
		return false
	}
	return claims.ExpiresAt.Before(time.Now())
}
