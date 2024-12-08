package middleware

import (
	"net/http"

	"github.com/shekshuev/shortener/internal/app/jwt"
	"github.com/shekshuev/shortener/internal/app/logger"
)

func RequestAuth(h http.Handler) http.Handler {
	log := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := jwt.GetAuthCookie(r)
		if err != nil {
			value, err := jwt.BuildJWTString()
			if err != nil {
				log.Log.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			cookie := &http.Cookie{
				Name:     jwt.CookieName,
				Value:    value,
				Path:     "/",
				HttpOnly: true,
			}
			r.AddCookie(cookie)
			http.SetCookie(w, cookie)
			jwtToken = value
			h.ServeHTTP(w, r)
			return
		}
		if jwt.IsTokenExpired(jwtToken) {
			http.SetCookie(w, &http.Cookie{
				Name:     jwt.CookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			w.WriteHeader(http.StatusUnauthorized)
			h.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}
