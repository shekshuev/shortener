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
	enableHTTPS := "true"
	cert := "cert"
	key := "key"
	os.Setenv("SERVER_ADDRESS", serverAddress)
	os.Setenv("BASE_URL", baseURL)
	os.Setenv("FILE_STORAGE_PATH", fileStoragePath)
	os.Setenv("DATABASE_DSN", databaseDSN)
	os.Setenv("ENABLE_HTTPS", enableHTTPS)
	os.Setenv("TLS_CERT", cert)
	os.Setenv("TLS_KEY", key)
	defer os.Unsetenv("SERVER_ADDRESS")
	defer os.Unsetenv("BASE_URL")
	defer os.Unsetenv("FILE_STORAGE_PATH")
	defer os.Unsetenv("DATABASE_DSN")
	defer os.Unsetenv("ENABLE_HTTPS")
	defer os.Unsetenv("TLS_CERT")
	defer os.Unsetenv("TLS_KEY")
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, baseURL)
	assert.Equal(t, cfg.ServerAddress, serverAddress)
	assert.Equal(t, cfg.FileStoragePath, fileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, databaseDSN)
	assert.Equal(t, cfg.EnableHTTPS, true)
	assert.Equal(t, cfg.CertFile, cert)
	assert.Equal(t, cfg.KeyFile, key)
}

func TestGetConfig_FlagPriority(t *testing.T) {
	serverAddress := "localhost:3000"
	baseURL := "http://localhost:3000"
	fileStoragePath := "./test.txt"
	databaseDSN := "host=test port=5432 user=test password=test dbname=test sslmode=disable"
	enableHTTPS := true
	cert := "cert"
	key := "key"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-a", serverAddress, "-b", baseURL, "-f", fileStoragePath, "-d", databaseDSN, "-s", "-C", cert, "-k", key}
	defer func() { os.Args = os.Args[:1] }()
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, baseURL)
	assert.Equal(t, cfg.ServerAddress, serverAddress)
	assert.Equal(t, cfg.FileStoragePath, fileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, databaseDSN)
	assert.Equal(t, cfg.EnableHTTPS, enableHTTPS)
	assert.Equal(t, cfg.CertFile, cert)
	assert.Equal(t, cfg.KeyFile, key)
}

func TestGetConfig_DefaultPriority(t *testing.T) {
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("ENABLE_HTTPS")
	os.Unsetenv("TLS_CERT")
	os.Unsetenv("TLS_KEY")
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd"}
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, cfg.DefaultBaseURL)
	assert.Equal(t, cfg.ServerAddress, cfg.DefaultServerAddress)
	assert.Equal(t, cfg.FileStoragePath, cfg.DefaultFileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, cfg.DefaultDatabaseDSN)
	assert.Equal(t, cfg.EnableHTTPS, cfg.DefaultEnableHTTPS)
	assert.Equal(t, cfg.CertFile, cfg.CertFile)
	assert.Equal(t, cfg.KeyFile, cfg.KeyFile)
}
