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
	subnet := "10.0.0.0/24"
	os.Setenv("SERVER_ADDRESS", serverAddress)
	os.Setenv("BASE_URL", baseURL)
	os.Setenv("FILE_STORAGE_PATH", fileStoragePath)
	os.Setenv("DATABASE_DSN", databaseDSN)
	os.Setenv("ENABLE_HTTPS", enableHTTPS)
	os.Setenv("TLS_CERT", cert)
	os.Setenv("TLS_KEY", key)
	os.Setenv("TRUSTED_SUBNET", subnet)
	defer os.Unsetenv("SERVER_ADDRESS")
	defer os.Unsetenv("BASE_URL")
	defer os.Unsetenv("FILE_STORAGE_PATH")
	defer os.Unsetenv("DATABASE_DSN")
	defer os.Unsetenv("ENABLE_HTTPS")
	defer os.Unsetenv("TLS_CERT")
	defer os.Unsetenv("TLS_KEY")
	defer os.Unsetenv("TRUSTED_SUBNET")
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, baseURL)
	assert.Equal(t, cfg.ServerAddress, serverAddress)
	assert.Equal(t, cfg.FileStoragePath, fileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, databaseDSN)
	assert.Equal(t, cfg.EnableHTTPS, true)
	assert.Equal(t, cfg.CertFile, cert)
	assert.Equal(t, cfg.KeyFile, key)
	assert.Equal(t, cfg.TrustedSubnet, subnet)
}

func TestGetConfig_FlagPriority(t *testing.T) {
	serverAddress := "localhost:3000"
	baseURL := "http://localhost:3000"
	fileStoragePath := "./test.txt"
	databaseDSN := "host=test port=5432 user=test password=test dbname=test sslmode=disable"
	enableHTTPS := true
	cert := "cert"
	key := "key"
	subnet := "10.0.0.0/24"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-a", serverAddress, "-b", baseURL, "-f", fileStoragePath, "-d", databaseDSN, "-s", "-cert", cert, "-key", key, "-t", subnet}
	defer func() { os.Args = os.Args[:1] }()
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, baseURL)
	assert.Equal(t, cfg.ServerAddress, serverAddress)
	assert.Equal(t, cfg.FileStoragePath, fileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, databaseDSN)
	assert.Equal(t, cfg.EnableHTTPS, enableHTTPS)
	assert.Equal(t, cfg.CertFile, cert)
	assert.Equal(t, cfg.KeyFile, key)
	assert.Equal(t, cfg.TrustedSubnet, subnet)

}

func TestGetConfig_DefaultPriority(t *testing.T) {
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("ENABLE_HTTPS")
	os.Unsetenv("TLS_CERT")
	os.Unsetenv("TLS_KEY")
	os.Unsetenv("TRUSTED_SUBNET")
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd"}
	cfg := GetConfig()
	assert.Equal(t, cfg.BaseURL, cfg.DefaultBaseURL)
	assert.Equal(t, cfg.ServerAddress, cfg.DefaultServerAddress)
	assert.Equal(t, cfg.FileStoragePath, cfg.DefaultFileStoragePath)
	assert.Equal(t, cfg.DatabaseDSN, cfg.DefaultDatabaseDSN)
	assert.Equal(t, cfg.EnableHTTPS, cfg.DefaultEnableHTTPS)
	assert.Equal(t, cfg.CertFile, cfg.DefaultCertFile)
	assert.Equal(t, cfg.KeyFile, cfg.DefaultKeyFile)
	assert.Equal(t, cfg.TrustedSubnet, cfg.DefaultTrustedSubnet)
}

func TestGetConfig_JSONPriority(t *testing.T) {
	// Создаём временный файл с JSON-конфигурацией
	tmpFile, err := os.CreateTemp("", "config*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	jsonContent := `{
		"server_address": "localhost:9999",
		"base_url": "http://json",
		"file_storage_path": "json_path.txt",
		"database_dsn": "json_dsn",
		"enable_https": true,
		"cert_file": "json_cert.pem",
		"key_file": "json_key.pem",
		"trusted_subnet": "10.0.0.0/24"
	}`
	_, err = tmpFile.WriteString(jsonContent)
	assert.NoError(t, err)
	tmpFile.Close()

	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("ENABLE_HTTPS")
	os.Unsetenv("TLS_CERT")
	os.Unsetenv("TLS_KEY")
	os.Unsetenv("TRUSTED_SUBNET")

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-c", tmpFile.Name()}

	cfg := GetConfig()

	assert.Equal(t, cfg.ServerAddress, "localhost:9999")
	assert.Equal(t, cfg.BaseURL, "http://json")
	assert.Equal(t, cfg.FileStoragePath, "json_path.txt")
	assert.Equal(t, cfg.DatabaseDSN, "json_dsn")
	assert.Equal(t, cfg.EnableHTTPS, true)
	assert.Equal(t, cfg.CertFile, "json_cert.pem")
	assert.Equal(t, cfg.KeyFile, "json_key.pem")
	assert.Equal(t, cfg.TrustedSubnet, "10.0.0.0/24")
}
