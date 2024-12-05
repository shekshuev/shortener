package store

import (
	"sync"

	"go.uber.org/zap"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/logger"
	"github.com/shekshuev/shortener/internal/app/models"
)

type UserURL struct {
	UserID string
	URL    string
}

type MemoryURLStore struct {
	mx   sync.RWMutex
	urls map[string]UserURL
	cfg  *config.Config
}

func NewMemoryURLStore(cfg *config.Config) *MemoryURLStore {
	store := &MemoryURLStore{urls: make(map[string]UserURL), cfg: cfg}
	log := logger.NewLogger()
	err := store.LoadSnapshot()
	if err != nil {
		log.Log.Error("Error loading snapshot", zap.Error(err))
	}
	return store
}

func (s *MemoryURLStore) SetURL(key, value, userID string) (string, error) {
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
	if len(userID) == 0 {
		return "", ErrEmptyUserID
	}
	s.urls[key] = UserURL{UserID: userID, URL: value}
	return value, nil
}

func (s *MemoryURLStore) SetBatchURL(createDTO []models.BatchShortURLCreateDTO, userID string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.urls == nil {
		return ErrNotInitialized
	}
	if len(userID) == 0 {
		return ErrEmptyUserID
	}
	if createDTO == nil {
		return ErrEmptyValue
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
		s.urls[dto.ShortURL] = UserURL{UserID: userID, URL: dto.OriginalURL}
	}
	return nil
}

func (s *MemoryURLStore) GetURL(key, userID string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if s.urls == nil {
		return "", ErrNotInitialized
	}
	value, exists := s.urls[key]
	if !exists {
		return "", ErrNotFound
	}
	if value.UserID != userID {
		return "", ErrNotFound
	}
	return value.URL, nil
}

func (s *MemoryURLStore) Close() error {
	return s.CreateSnapshot()
}
