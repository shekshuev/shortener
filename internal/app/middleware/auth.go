package middleware

import (
	"net/http"

	"github.com/shekshuev/shortener/internal/app/jwt"
	"github.com/shekshuev/shortener/internal/app/logger"
)

const cookieName = "token"

func RequestAuth(h http.Handler) http.Handler {
	log := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := r.Cookies()
		for _, value := range cookie {
			if value.Name != cookieName {
				continue
			}
			if jwt.IsTokenExpired(value.Value) {
				http.SetCookie(w, &http.Cookie{
					Name:     cookieName,
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
				h.ServeHTTP(w, r)
				return
			}
		}
		value, err := jwt.BuildJWTString()
		if err != nil {
			log.Log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    value,
			Path:     "/",
			HttpOnly: true,
		})
		h.ServeHTTP(w, r)
	})
}
