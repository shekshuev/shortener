package store

import (
	"sync"

	"go.uber.org/zap"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/logger"
	"github.com/shekshuev/shortener/internal/app/models"
)

type MemoryURLStore struct {
	mx   sync.RWMutex
	urls map[string]string
	cfg  *config.Config
}

func NewMemoryURLStore(cfg *config.Config) *MemoryURLStore {
	store := &MemoryURLStore{urls: make(map[string]string), cfg: cfg}
	log := logger.NewLogger()
	err := store.LoadSnapshot()
	if err != nil {
		log.Log.Error("Error loading snapshot", zap.Error(err))
	}
	return store
}

func (s *MemoryURLStore) SetURL(key, value string) (string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.urls == nil {
		return "", ErrNotInitialized
	}
	if len(key) == 0 {
		return "", ErrEmptyKey
	}
	if len(value) == 0 {
		return "", ErrEmptyValue
	}
	s.urls[key] = value
	return value, nil
}

func (s *MemoryURLStore) SetBatchURL(createDTO []models.BatchShortURLCreateDTO) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.urls == nil {
		return ErrNotInitialized
	}
	for _, dto := range createDTO {
		if len(dto.ShortURL) == 0 {
			return ErrEmptyKey
		}
		if len(dto.OriginalURL) == 0 {
			return ErrEmptyValue
		}
	}
	for _, dto := range createDTO {
		s.urls[dto.ShortURL] = dto.OriginalURL
	}
	return nil
}

func (s *MemoryURLStore) GetURL(key string) (string, error) {
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

func (s *MemoryURLStore) Close() error {
	return s.CreateSnapshot()
}
