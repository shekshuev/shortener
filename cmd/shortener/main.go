package main

import (
	"context"
	"fmt"
	"net"
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
	"google.golang.org/grpc"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	_ "net/http/pprof"

	"github.com/shekshuev/shortener/internal/app/grpcserver"
	"github.com/shekshuev/shortener/internal/app/proto"
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

	var trustedSubnet *net.IPNet
	if cfg.TrustedSubnet != "" {
		_, subnet, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			l.Log.Fatal("Invalid trusted subnet", zap.Error(err))
		}
		trustedSubnet = subnet
	}

	var urlStore store.URLStore
	if cfg.DatabaseDSN == cfg.DefaultDatabaseDSN {
		urlStore = store.NewMemoryURLStore(&cfg)
	} else {
		urlStore = store.NewPostgresURLStore(&cfg)
	}
	urlService := service.NewURLService(urlStore, &cfg)

	urlHandler := handler.NewURLHandler(urlService, trustedSubnet)
	httpServer := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: urlHandler.Router,
	}

	grpcSrv := grpc.NewServer()
	grpcHandler := grpcserver.NewServer(urlService)
	proto.RegisterURLShortenerServer(grpcSrv, grpcHandler)

	grpcListener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		l.Log.Fatal("Failed to listen for gRPC", zap.Error(err))
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		var err error
		if cfg.EnableHTTPS {
			l.Log.Info("Starting HTTPS server", zap.String("addr", cfg.ServerAddress))
			err = httpServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			l.Log.Info("Starting HTTP server", zap.String("addr", cfg.ServerAddress))
			err = httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			l.Log.Error("Error starting HTTP server", zap.Error(err))
		}
	}()

	go func() {
		l.Log.Info("Starting gRPC server", zap.String("addr", cfg.GRPCServerAddress))
		if err := grpcSrv.Serve(grpcListener); err != nil {
			l.Log.Error("Error starting gRPC server", zap.Error(err))
		}
	}()

	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	l.Log.Info("Servers started. Waiting for shutdown signal...")

	<-done
	l.Log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcSrv.GracefulStop()
	l.Log.Info("gRPC server shutdown gracefully")

	if err := httpServer.Shutdown(ctx); err != nil {
		l.Log.Error("HTTP server forced to shutdown", zap.Error(err))
	} else {
		l.Log.Info("HTTP server shutdown gracefully")
	}

	if err := urlStore.Close(); err != nil {
		l.Log.Error("Error closing store", zap.Error(err))
	} else {
		l.Log.Info("Store closed")
	}
}
