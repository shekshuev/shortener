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
	router.Use(middleware.RequestLogger)
	router.Use(middleware.GzipCompressor)
	h := &URLHandler{service: service, Router: router}
	router.Post("/", h.createURLHandler)
	router.Post("/api/shorten", h.createURLHandlerJSON)
	router.Get("/{shorted}", h.getURLHandler)
	return h
}

func (h *URLHandler) createURLHandler(w http.ResponseWriter, r *http.Request) {
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
}

func (h *URLHandler) createURLHandlerJSON(w http.ResponseWriter, r *http.Request) {
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
}

func (h *URLHandler) getURLHandler(w http.ResponseWriter, r *http.Request) {
	if longURL, err := h.service.GetLongURL(r.URL.Path[1:]); err == nil {
		http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}
