package logger

import (
	"sync"

	"go.uber.org/zap"
)

type Logger struct {
	Log *zap.Logger
}

var (
	instance *Logger
	once     sync.Once
)

func GetInstance() *Logger {
	once.Do(func() {
		instance = &Logger{Log: zap.NewNop()}
	})
	return instance
}

func (l *Logger) Initialize(level string) error {
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
