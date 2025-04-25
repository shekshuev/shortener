package jwt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestBuildJWTString(t *testing.T) {
	tokenStr, err := BuildJWTString()
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)
}

func TestGetAuthCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := &http.Cookie{
		Name:  CookieName,
		Value: "mytoken",
	}
	req.AddCookie(c)

	token, err := GetAuthCookie(req)
	assert.NoError(t, err)
	assert.Equal(t, "mytoken", token)
}

func TestGetAuthCookie_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := GetAuthCookie(req)
	assert.Error(t, err)
}

func TestGetUserID(t *testing.T) {
	tokenStr, err := BuildJWTString()
	assert.NoError(t, err)

	userID, err := GetUserID(tokenStr)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)
}

func TestIsTokenExpired_ValidToken(t *testing.T) {
	tokenStr, err := buildJWTWithCustomExp(time.Now().Add(time.Hour))
	assert.NoError(t, err)

	expired := IsTokenExpired(tokenStr)
	assert.False(t, expired)
}

func TestIsTokenExpired_ExpiredToken(t *testing.T) {
	exp := time.Now().Add(-1 * time.Hour)
	fmt.Println("exp =", exp)

	tokenStr, err := buildJWTWithCustomExp(exp)
	assert.NoError(t, err)

	expired := IsTokenExpired(tokenStr)
	assert.True(t, expired, "Expected token to be expired")
}

func TestIsTokenExpired_InvalidToken(t *testing.T) {
	expired := IsTokenExpired("invalid.token.string")
	assert.False(t, expired, "Should return false on parse error")
}

// 🔧 вспомогательная функция
func buildJWTWithCustomExp(exp time.Time) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		UserID: "test-user",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}
