package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/shekshuev/shortener/internal/app/jwt"
	"github.com/stretchr/testify/assert"
)

func buildJWTWithExpiration(exp time.Time) (string, error) {
	claims := jwt.Claims{
		RegisteredClaims: jwtv4.RegisteredClaims{
			ExpiresAt: jwtv4.NewNumericDate(exp),
		},
		UserID: uuid.New().String(),
	}
	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwt.SecretKey))
}

func TestRequestAuth_NoCookie_SetsNewJWT(t *testing.T) {
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware := RequestAuth(handler)
	middleware.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.True(t, handlerCalled)
	found := false
	for _, c := range resp.Cookies() {
		if c.Name == jwt.CookieName {
			found = true
			assert.NotEmpty(t, c.Value)
		}
	}
	assert.True(t, found, "JWT cookie should be set")
}

func TestRequestAuth_ExpiredCookie_DeletesCookieAndReturns401(t *testing.T) {
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	expiredToken, err := buildJWTWithExpiration(time.Now().Add(-time.Hour))
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  jwt.CookieName,
		Value: expiredToken,
	})
	rec := httptest.NewRecorder()

	middleware := RequestAuth(handler)
	middleware.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	found := false
	for _, c := range resp.Cookies() {
		if c.Name == jwt.CookieName && c.MaxAge == -1 {
			found = true
		}
	}
	assert.True(t, found, "Expired cookie should be cleared")
}

func TestRequestAuth_ValidCookie_PassesThrough(t *testing.T) {
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	validToken, err := buildJWTWithExpiration(time.Now().Add(time.Hour))
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  jwt.CookieName,
		Value: validToken,
	})
	rec := httptest.NewRecorder()

	middleware := RequestAuth(handler)
	middleware.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
