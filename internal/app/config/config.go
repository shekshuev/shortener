package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type envConfig struct {
	FlagRunAddr     string `env:"SERVER_ADDRESS"`
	BaseShorterAddr string `env:"BASE_URL"`
}

var (
	FlagRunAddr     string
	BaseShorterAddr string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&BaseShorterAddr, "b", "http://localhost:8080", "base url of shorter address")
	flag.Parse()
	parsEnv()
}

func SetConfig(runAddr, baseAddr string) {
	FlagRunAddr = runAddr
	BaseShorterAddr = baseAddr
	parsEnv()
}

func parsEnv() {
	var cfg envConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if len(cfg.BaseShorterAddr) > 0 {
		BaseShorterAddr = cfg.BaseShorterAddr
	}
	if len(cfg.FlagRunAddr) > 0 {
		BaseShorterAddr = cfg.FlagRunAddr
	}
}
