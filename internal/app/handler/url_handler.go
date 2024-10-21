package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shekshuev/shortener/internal/app/middleware"
	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/shekshuev/shortener/internal/app/service"
)

type URLHandler struct {
	service *service.URLService
	Router  *chi.Mux
}

func NewURLHandler(service *service.URLService) *URLHandler {
	router := chi.NewRouter()
	h := &URLHandler{service: service, Router: router}
	router.Post("/", middleware.RequestLogger(h.createURLHandler))
	router.Post("/api/shorten", middleware.RequestLogger(h.createURLHandlerJSON))
	router.Get("/{shorted}", middleware.RequestLogger(h.getURLHandler))
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
		shortURL, err := h.service.CreateShortURL(string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(shortURL))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *URLHandler) createURLHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var createDTO models.ShortURLCreateDTO
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(body, &createDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		shortURL, err := h.service.CreateShortURL(createDTO.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(shortURL))
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
