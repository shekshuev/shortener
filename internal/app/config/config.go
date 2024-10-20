package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/shekshuev/shortener/internal/app/logger"
	"go.uber.org/zap"
)

type Config struct {
	ServerAddress        string
	BaseURL              string
	DefaultServerAddress string
	DefaultBaseURL       string
}

type envConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func GetConfig() Config {

	var cfg Config
	cfg.DefaultServerAddress = "localhost:8080"
	cfg.DefaultBaseURL = "http://localhost:8080"
	parseFlags(&cfg)
	parsEnv(&cfg)
	return cfg
}

func parseFlags(cfg *Config) {
	if f := flag.Lookup("a"); f == nil {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.DefaultServerAddress, "address and port to run server")
	} else {
		cfg.ServerAddress = cfg.DefaultServerAddress
	}
	if f := flag.Lookup("b"); f == nil {
		flag.StringVar(&cfg.BaseURL, "b", cfg.DefaultBaseURL, "base URL of shorter address")
	} else {
		cfg.BaseURL = cfg.DefaultBaseURL
	}
	flag.Parse()
	parsEnv(cfg)
}

func parsEnv(cfg *Config) {
	var envCfg envConfig
	err := env.Parse(&envCfg)
	if err != nil {
		logger.Log.Error("Error starting server", zap.Error(err))
	}
	if len(envCfg.BaseURL) > 0 {
		cfg.BaseURL = envCfg.BaseURL
	}
	if len(envCfg.ServerAddress) > 0 {
		cfg.ServerAddress = envCfg.ServerAddress
	}
}
