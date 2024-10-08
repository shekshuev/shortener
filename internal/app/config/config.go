package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	FlagRunAddr     string
	BaseShorterAddr string
}

type envConfig struct {
	FlagRunAddr     string `env:"SERVER_ADDRESS"`
	BaseShorterAddr string `env:"BASE_URL"`
}

func GetConfig() Config {
	var cfg Config
	parseFlags(&cfg)
	parsEnv(&cfg)
	return cfg
}

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.BaseShorterAddr, "b", "http://localhost:8080", "base url of shorter address")
	flag.Parse()
	parsEnv(cfg)
}

func parsEnv(cfg *Config) {
	var envCfg envConfig
	err := env.Parse(&envCfg)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	if len(envCfg.BaseShorterAddr) > 0 {
		cfg.BaseShorterAddr = envCfg.BaseShorterAddr
	}
	if len(envCfg.FlagRunAddr) > 0 {
		cfg.FlagRunAddr = envCfg.FlagRunAddr
	}
}
