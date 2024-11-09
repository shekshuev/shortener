package store

import (
	"fmt"
	"sync"

	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/logger"
)

type URLStore struct {
	mx   sync.RWMutex
	urls map[string]string
	cfg  *config.Config
	db   *sql.DB
}

var ErrEmptyKey = fmt.Errorf("key cannot be empty")
var ErrEmptyValue = fmt.Errorf("value cannot be empty")
var ErrNotFound = fmt.Errorf("not found")
var ErrNotInitialized = fmt.Errorf("store not initialized")

func NewURLStore(cfg *config.Config) *URLStore {
	log := logger.NewLogger()
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Log.Fatal("Error connecting to database", zap.Error(err))
	}
	store := &URLStore{urls: make(map[string]string), cfg: cfg, db: db}
	err = store.LoadSnapshot()
	if err != nil {
		log.Log.Error("Error loading snapshot", zap.Error(err))
	}
	return store
}

func (s *URLStore) SetURL(key, value string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.urls == nil {
		return ErrNotInitialized
	}
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if len(value) == 0 {
		return ErrEmptyValue
	}
	s.urls[key] = value
	return nil
}

func (s *URLStore) GetURL(key string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if s.urls == nil {
		return "", ErrNotInitialized
	}
	value, exists := s.urls[key]
	if !exists {
		return "", ErrNotFound
	}
	return value, nil
}

func (s *URLStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *URLStore) CheckDBConnection() error {
	return s.db.Ping()
}
