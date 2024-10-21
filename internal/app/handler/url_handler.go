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

type Middleware func(http.HandlerFunc) http.HandlerFunc

func conveyor(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func NewURLHandler(service *service.URLService) *URLHandler {
	router := chi.NewRouter()
	h := &URLHandler{service: service, Router: router}
	router.Post("/", conveyor(h.createURLHandler, middleware.RequestLogger, middleware.GzipCompressor))
	router.Post("/api/shorten", conveyor(h.createURLHandlerJSON, middleware.RequestLogger, middleware.GzipCompressor))
	router.Get("/{shorted}", conveyor(h.getURLHandler, middleware.RequestLogger, middleware.GzipCompressor))
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
		w.Header().Set("Content-Type", "application/json")
		shortURL, err := h.service.CreateShortURL(createDTO.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusCreated)
		readDTO := models.ShortURLReadDTO{Result: shortURL}
		resp, err := json.Marshal(readDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		_, err = w.Write(resp)
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
