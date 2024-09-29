package main

import (
	"fmt"
	"github.com/shekshuev/shortener/internal/utils"
	"io"
	"net/http"
)

type URLShortener struct {
	urls map[string]string
}

func (us *URLShortener) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		shorted := utils.Shorten(string(body))
		us.urls[shorted] = string(body)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%s/%s", r.Host, shorted)))
	} else if r.Method == http.MethodGet {
		shorted := r.URL.Path[1:]
		if url, ok := us.urls[shorted]; ok {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	mux := http.NewServeMux()
	us := &URLShortener{
		urls: make(map[string]string),
	}

	mux.HandleFunc("/", us.handler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
