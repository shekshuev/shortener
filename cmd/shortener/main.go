package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/handler"
	"github.com/shekshuev/shortener/internal/app/logger"
	"github.com/shekshuev/shortener/internal/app/service"
	"github.com/shekshuev/shortener/internal/app/store"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	_ "net/http/pprof"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {
	printBuildInfo()

	l := logger.NewLogger()
	cfg := config.GetConfig()
	var urlStore store.URLStore = nil
	if cfg.DatabaseDSN == cfg.DefaultDatabaseDSN {
		urlStore = store.NewMemoryURLStore(&cfg)
	} else {
		urlStore = store.NewPostgresURLStore(&cfg)
	}
	urlService := service.NewURLService(urlStore, &cfg)
	urlHandler := handler.NewURLHandler(urlService)
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: urlHandler.Router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		var err error
		if cfg.EnableHTTPS {
			l.Log.Info("Starting HTTPS server", zap.String("addr", cfg.ServerAddress))
			err = server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			l.Log.Info("Starting HTTP server", zap.String("addr", cfg.ServerAddress))
			err = server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			l.Log.Error("Error starting server", zap.Error(err))
		}
	}()
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	l.Log.Info("Server started. Waiting for shutdown signal...")

	<-done
	l.Log.Info("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := urlStore.Close(); err != nil {
		l.Log.Error("Error closing store", zap.Error(err))
	} else {
		l.Log.Info("Store closed")
	}
	if err := server.Shutdown(ctx); err != nil {
		l.Log.Error("Server forced to shutdown", zap.Error(err))
	} else {
		l.Log.Info("Server shutdown gracefully")
	}
}
