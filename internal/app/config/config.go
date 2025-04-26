package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/shekshuev/shortener/internal/app/logger"
	"go.uber.org/zap"
)

// Config содержит настройки приложения, включая параметры сервера, базы данных и файлового хранилища.
type Config struct {
	ServerAddress          string // Адрес и порт, на котором запускается сервер.
	BaseURL                string // Базовый URL для сокращённых ссылок.
	FileStoragePath        string // Путь к файлу для хранения сокращённых URL.
	DatabaseDSN            string // Строка подключения к базе данных.
	EnableHTTPS            bool   // Включить HTTPS.
	CertFile               string // Путь к файлу с сертификатом.
	KeyFile                string // путь к файлу с ключом.
	TrustedSubnet          string // Доверенная подсеть в CIDR-формате.
	DefaultServerAddress   string // Значение по умолчанию для ServerAddress.
	DefaultBaseURL         string // Значение по умолчанию для BaseURL.
	DefaultFileStoragePath string // Значение по умолчанию для FileStoragePath.
	DefaultDatabaseDSN     string // Значение по умолчанию для DatabaseDSN.
	DefaultEnableHTTPS     bool   // Значение по умолчанию для EnableHTTPS.
	DefaultCertFile        string // Значение по умолчанию для CertFile.
	DefaultKeyFile         string // Значение по умолчанию для KeyFile.
	DefaultTrustedSubnet   string // Значение по умолчанию для TrustedSubnet (пустая строка).
}

type envConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	EnableHTTPS     string `env:"ENABLE_HTTPS"`
	CertFile        string `env:"TLS_CERT"`
	KeyFile         string `env:"TLS_KEY"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`
}

type jsonConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	CertFile        string `json:"cert_file"`
	KeyFile         string `json:"key_file"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

// GetConfig возвращает экземпляр конфига
func GetConfig() Config {
	var cfg Config
	cfg.DefaultServerAddress = "localhost:8080"
	cfg.DefaultBaseURL = "http://localhost:8080"
	cfg.DefaultFileStoragePath = "./storage.txt"
	cfg.DefaultDatabaseDSN = ""
	cfg.DefaultEnableHTTPS = false
	cfg.DefaultCertFile = ""
	cfg.DefaultKeyFile = ""
	cfg.DefaultTrustedSubnet = ""
	parseFlags(&cfg)
	parsEnv(&cfg)
	return cfg
}

func parseFlags(cfg *Config) {
	var configPath string
	if f := flag.Lookup("c"); f == nil {
		flag.StringVar(&configPath, "c", "", "path to JSON config file")
	} else if f := flag.Lookup("config"); f == nil {
		flag.StringVar(&configPath, "config", "", "path to JSON config file")
	}
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
	if f := flag.Lookup("f"); f == nil {
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.DefaultFileStoragePath, "file storage path")
	} else {
		cfg.FileStoragePath = cfg.DefaultFileStoragePath
	}
	if f := flag.Lookup("d"); f == nil {
		flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DefaultDatabaseDSN, "database connection string")
	} else {
		cfg.DatabaseDSN = cfg.DefaultDatabaseDSN
	}
	if f := flag.Lookup("s"); f == nil {
		flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.DefaultEnableHTTPS, "enable HTTPS")
	} else {
		cfg.EnableHTTPS = cfg.DefaultEnableHTTPS
	}

	if f := flag.Lookup("cert"); f == nil {
		flag.StringVar(&cfg.CertFile, "cert", cfg.DefaultCertFile, "cert file")
	} else {
		cfg.CertFile = cfg.DefaultCertFile
	}
	if f := flag.Lookup("key"); f == nil {
		flag.StringVar(&cfg.KeyFile, "key", cfg.DefaultKeyFile, "key file")
	} else {
		cfg.KeyFile = cfg.DefaultKeyFile
	}
	if f := flag.Lookup("t"); f == nil {
		flag.StringVar(&cfg.TrustedSubnet, "t", cfg.DefaultTrustedSubnet, "trusted subnet CIDR")
	} else {
		cfg.TrustedSubnet = cfg.DefaultTrustedSubnet
	}
	flag.Parse()
	parseJSON(configPath, cfg)
	parsEnv(cfg)
}

func parsEnv(cfg *Config) {
	l := logger.NewLogger()
	var envCfg envConfig
	err := env.Parse(&envCfg)
	if err != nil {
		l.Log.Error("Error starting server", zap.Error(err))
	}
	if len(envCfg.BaseURL) > 0 {
		cfg.BaseURL = envCfg.BaseURL
	}
	if len(envCfg.ServerAddress) > 0 {
		cfg.ServerAddress = envCfg.ServerAddress
	}
	if len(envCfg.FileStoragePath) > 0 {
		cfg.FileStoragePath = envCfg.FileStoragePath
	}
	if len(envCfg.DatabaseDSN) > 0 {
		cfg.DatabaseDSN = envCfg.DatabaseDSN
	}
	if envCfg.EnableHTTPS == "true" || envCfg.EnableHTTPS == "1" {
		cfg.EnableHTTPS = true
	}
	if len(envCfg.CertFile) > 0 {
		cfg.CertFile = envCfg.CertFile
	}
	if len(envCfg.KeyFile) > 0 {
		cfg.KeyFile = envCfg.KeyFile
	}
	if len(envCfg.TrustedSubnet) > 0 {
		cfg.TrustedSubnet = envCfg.TrustedSubnet
	}
}

func parseJSON(path string, cfg *Config) {
	if path == "" {
		path = os.Getenv("CONFIG")
	}
	if path == "" {
		return
	}

	file, err := os.Open(path)
	if err != nil {
		logger.NewLogger().Log.Warn("Could not open config file", zap.Error(err))
		return
	}
	defer file.Close()

	var jCfg jsonConfig
	if err := json.NewDecoder(file).Decode(&jCfg); err != nil {
		logger.NewLogger().Log.Warn("Could not decode config JSON", zap.Error(err))
		return
	}

	if cfg.ServerAddress == cfg.DefaultServerAddress && jCfg.ServerAddress != "" {
		cfg.ServerAddress = jCfg.ServerAddress
	}
	if cfg.BaseURL == cfg.DefaultBaseURL && jCfg.BaseURL != "" {
		cfg.BaseURL = jCfg.BaseURL
	}
	if cfg.FileStoragePath == cfg.DefaultFileStoragePath && jCfg.FileStoragePath != "" {
		cfg.FileStoragePath = jCfg.FileStoragePath
	}
	if cfg.DatabaseDSN == cfg.DefaultDatabaseDSN && jCfg.DatabaseDSN != "" {
		cfg.DatabaseDSN = jCfg.DatabaseDSN
	}
	if cfg.EnableHTTPS == cfg.DefaultEnableHTTPS {
		cfg.EnableHTTPS = jCfg.EnableHTTPS
	}
	if cfg.CertFile == cfg.DefaultCertFile && jCfg.CertFile != "" {
		cfg.CertFile = jCfg.CertFile
	}
	if cfg.KeyFile == cfg.DefaultKeyFile && jCfg.KeyFile != "" {
		cfg.KeyFile = jCfg.KeyFile
	}
	if cfg.TrustedSubnet == cfg.DefaultTrustedSubnet && jCfg.TrustedSubnet != "" {
		cfg.TrustedSubnet = jCfg.TrustedSubnet
	}
}
