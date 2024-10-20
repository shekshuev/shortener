package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig_EnvPriority(t *testing.T) {
	serverAddress := "localhost:3000"
	baseURL := "http://localhost:3000"
	os.Setenv("SERVER_ADDRESS", serverAddress)
	os.Setenv("BASE_URL", baseURL)
	defer os.Unsetenv("SERVER_ADDRESS")
	defer os.Unsetenv("BASE_URL")
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, baseURL)
	assert.Equal(t, cfg.ServerAddress, serverAddress)
}

func TestGetConfig_FlagPriority(t *testing.T) {
	serverAddress := "localhost:3000"
	baseURL := "http://localhost:3000"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-a", serverAddress, "-b", baseURL}
	defer func() { os.Args = os.Args[:1] }()
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, baseURL)
	assert.Equal(t, cfg.ServerAddress, serverAddress)
}

func TestGetConfig_DefaultPriority(t *testing.T) {
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd"}
	cfg := GetConfig()
	if cfg.ServerAddress != "localhost:8080" {
		t.Errorf("Expected SERVER_ADDRESS to be '%s', got '%s'", cfg.DefaultServerAddress, cfg.ServerAddress)
	}
	if cfg.BaseURL != "http://localhost:8080" {
		t.Errorf("Expected BASE_URL to be '%s', got '%s'", cfg.DefaultBaseURL, cfg.BaseURL)
	}
}
