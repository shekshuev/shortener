package main

import (
	"fmt"
	"github.com/shekshuev/shortener/internal/utils"
	"io"
	"net/http"
)

var urls = make(map[string]string)

func create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		shorted := utils.Shorten(string(body))
		urls[shorted] = string(body)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, shorted)))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if url, ok := urls[r.URL.Path[1:]]; ok {
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
	mux.HandleFunc("/", create)
	mux.HandleFunc("/{shorted}", get)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
