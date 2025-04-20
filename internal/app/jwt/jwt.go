package jwt

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// CookieName - имя куки, в котором хранится токен.
const CookieName = "token"

// Claims - структура, представляющая собой полезную нагрузку JWT-токена.
type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

// TokenExp - время жизни токена.
const TokenExp = time.Hour * 3

// SecretKey - секретный ключ для подписи токена.
const SecretKey = "supersecretkey"

// GetAuthCookie извлекает значение куки с токеном из HTTP-запроса.
func GetAuthCookie(req *http.Request) (string, error) {
	cookie, err := req.Cookie(CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// BuildJWTString создаёт новый JWT-токен и возвращает его строковое представление.
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

// fromString парсит строку токена и извлекает из него данные.
func fromString(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// func fromString(tokenString string) (*Claims, error) {
// 	claims := &Claims{}
// 	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(SecretKey), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return claims, nil
// }

// GetUserID извлекает UserID из токена.
func GetUserID(tokenString string) (string, error) {
	claims, err := fromString(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// IsTokenExpired проверяет, истёк ли срок действия токена.
func IsTokenExpired(tokenString string) bool {
	claims, err := fromString(tokenString)
	if err != nil {
		return false
	}
	if claims.ExpiresAt == nil {
		return false
	}
	return claims.ExpiresAt.Before(time.Now())
}
