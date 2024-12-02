package middleware

import (
	"net/http"

	"github.com/shekshuev/shortener/internal/app/jwt"
	"github.com/shekshuev/shortener/internal/app/logger"
)

func RequestAuth(h http.Handler) http.Handler {
	log := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := jwt.GetAuthCookie(r)
		if err != nil {
			value, err := jwt.BuildJWTString()
			if err != nil {
				log.Log.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     jwt.CookieName,
				Value:    value,
				Path:     "/",
				HttpOnly: true,
			})
			cookie = value
		}
		if jwt.IsTokenExpired(cookie) {
			http.SetCookie(w, &http.Cookie{
				Name:     jwt.CookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
		}
		h.ServeHTTP(w, r)
	})
}
