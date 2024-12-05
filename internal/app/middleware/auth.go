package middleware

import (
	"net/http"

	"github.com/shekshuev/shortener/internal/app/jwt"
	"github.com/shekshuev/shortener/internal/app/logger"
)

func RequestAuth(h http.Handler) http.Handler {
	log := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			h.ServeHTTP(w, r)
			return
		}
		jwtToken, err := jwt.GetAuthCookie(r)

		if err != nil && r.Method == http.MethodPost {
			value, err := jwt.BuildJWTString()
			if err != nil {
				log.Log.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
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
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
