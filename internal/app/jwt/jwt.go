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

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: uuid.New().String(),
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func fromString(tokenString string) *Claims {
	claims := &Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	return claims
}

func GetUserID(tokenString string) string {
	claims := fromString(tokenString)
	return claims.UserID
}

func IsTokenExpired(tokenString string) bool {
	claims := fromString(tokenString)
	if claims.ExpiresAt == nil {
		return false
	}
	return claims.ExpiresAt.Before(time.Now())
}
