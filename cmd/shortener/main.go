package main

import (
	"log"
	"net/http"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/handler"
	"github.com/shekshuev/shortener/internal/app/logger"
	"go.uber.org/zap"

	"github.com/shekshuev/shortener/internal/app/service"
	"github.com/shekshuev/shortener/internal/app/store"
)

func main() {
	if err := logger.Initialize("info"); err != nil {
		log.Fatalf("Error initialize zap logger: %v", err)
	}
	cfg := config.GetConfig()
	urlStore := store.NewURLStore()
	urlService := service.NewURLService(urlStore, &cfg)
	urlHandler := handler.NewURLHandler(urlService)
	if err := http.ListenAndServe(cfg.ServerAddress, urlHandler.Router); err != nil {
		logger.Log.Error("Error starting server", zap.Error(err))
	}
}
