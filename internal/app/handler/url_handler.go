package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shekshuev/shortener/internal/app/service"
)

type URLHandler struct {
	service *service.URLService
	Router  *chi.Mux
}

func NewURLHandler(service *service.URLService) *URLHandler {
	router := chi.NewRouter()
	h := &URLHandler{service: service, Router: router}
	router.Post("/", h.createURLHandler)
	router.Get("/{shorted}", h.getURLHandler)
	return h
}

func (h *URLHandler) createURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		shorted, err := h.service.CreateShortURL(string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, shorted)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *URLHandler) getURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if longURL, err := h.service.GetLongURL(r.URL.Path[1:]); err == nil {
			http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}
