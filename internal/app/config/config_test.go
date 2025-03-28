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
	fileStoragePath := "./test.txt"
	databaseDSN := "host=test port=5432 user=test password=test dbname=test sslmode=disable"
	os.Setenv("SERVER_ADDRESS", serverAddress)
	os.Setenv("BASE_URL", baseURL)
	os.Setenv("FILE_STORAGE_PATH", fileStoragePath)
	os.Setenv("DATABASE_DSN", databaseDSN)
	defer os.Unsetenv("SERVER_ADDRESS")
	defer os.Unsetenv("BASE_URL")
	defer os.Unsetenv("FILE_STORAGE_PATH")
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, baseURL)
	assert.Equal(t, cfg.ServerAddress, serverAddress)
	assert.Equal(t, cfg.FileStoragePath, fileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, databaseDSN)
}

func TestGetConfig_FlagPriority(t *testing.T) {
	serverAddress := "localhost:3000"
	baseURL := "http://localhost:3000"
	fileStoragePath := "./test.txt"
	databaseDSN := "host=test port=5432 user=test password=test dbname=test sslmode=disable"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-a", serverAddress, "-b", baseURL, "-f", fileStoragePath, "-d", databaseDSN}
	defer func() { os.Args = os.Args[:1] }()
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, baseURL)
	assert.Equal(t, cfg.ServerAddress, serverAddress)
	assert.Equal(t, cfg.FileStoragePath, fileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, databaseDSN)
}

func TestGetConfig_DefaultPriority(t *testing.T) {
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("DATABASE_DSN")
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd"}
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, cfg.DefaultBaseURL)
	assert.Equal(t, cfg.ServerAddress, cfg.DefaultServerAddress)
	assert.Equal(t, cfg.FileStoragePath, cfg.DefaultFileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, cfg.DefaultDatabaseDSN)
}
