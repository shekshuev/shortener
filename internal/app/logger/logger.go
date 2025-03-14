package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

// Logger - структура для логирования с использованием zap.
type Logger struct {
	Log *zap.Logger
}

var (
	instance *Logger
	once     sync.Once
)

// NewLogger создаёт и возвращает экземпляр логгера (Singleton).
func NewLogger() *Logger {
	once.Do(func() {
		instance = &Logger{Log: zap.NewNop()}
		if err := instance.initialize("info"); err != nil {
			log.Fatalf("Error initializing zap logger: %v", err)
		}
	})
	return instance
}

// initialize инициализирует zap-логгер с указанным уровнем логирования.
func (l *Logger) initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	l.Log = zl
	return nil
}
