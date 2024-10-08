package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
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
	f := flag.FlagSet{}
	f.StringVar(&cfg.ServerAddress, "a", cfg.DefaultServerAddress, "address and port to run server")
	f.StringVar(&cfg.BaseURL, "b", cfg.DefaultBaseURL, "base url of shorter address")
	err := f.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
	parsEnv(cfg)
}

func parsEnv(cfg *Config) {
	var envCfg envConfig
	err := env.Parse(&envCfg)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	if len(envCfg.BaseURL) > 0 {
		cfg.BaseURL = envCfg.BaseURL
	}
	if len(envCfg.ServerAddress) > 0 {
		cfg.ServerAddress = envCfg.ServerAddress
	}
}
