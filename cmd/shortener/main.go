package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/utils"
)

type Shortener struct {
	urls map[string]string
}

func (s *Shortener) createUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		shorted := utils.Shorten(string(body))
		s.urls[shorted] = string(body)
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, shorted)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *Shortener) getUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if url, ok := s.urls[r.URL.Path[1:]]; ok {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func main() {
	s := &Shortener{urls: make(map[string]string)}
	cfg := config.GetConfig()
	r := chi.NewRouter()
	r.Post("/", s.createUrlHandler)
	r.Get("/{shorted}", s.getUrlHandler)
	if err := http.ListenAndServe(cfg.FlagRunAddr, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
