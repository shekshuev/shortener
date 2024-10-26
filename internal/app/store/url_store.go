package store

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/logger"
)

type URLStore struct {
	mx   sync.RWMutex
	urls map[string]string
	cfg  *config.Config
}

var ErrEmptyKey = fmt.Errorf("key cannot be empty")
var ErrEmptyValue = fmt.Errorf("value cannot be empty")
var ErrNotFound = fmt.Errorf("not found")

func NewURLStore(cfg *config.Config) *URLStore {
	store := &URLStore{urls: make(map[string]string), cfg: cfg}
	log := logger.GetInstance()
	err := store.LoadSnapshot()
	if err != nil {
		log.Log.Error("Error loading snapshot", zap.Error(err))
	}
	return store
}

func (s *URLStore) SetURL(key, value string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if len(value) == 0 {
		return ErrEmptyValue
	}
	s.urls[key] = value
	return s.CreateSnapshot()
}

func (s *URLStore) GetURL(key string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	value, exists := s.urls[key]
	if !exists {
		return "", ErrNotFound
	}
	return value, nil
}
