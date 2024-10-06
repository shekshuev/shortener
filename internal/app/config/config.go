package config

import (
	"flag"
)

var (
	FlagRunAddr     string
	BaseShorterAddr string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&BaseShorterAddr, "b", "http://localhost:8080", "base url of shorter address")
	flag.Parse()
}

func SetConfig(runAddr, baseAddr string) {
	FlagRunAddr = runAddr
	BaseShorterAddr = baseAddr
}
