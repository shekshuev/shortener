package store

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/logger"
	"github.com/shekshuev/shortener/internal/app/models"
)

// UserURL представляет структуру для хранения информации о сокращённом URL.
type UserURL struct {
	UserID    string
	URL       string
	IsDeleted bool
}

// MemoryURLStore - хранилище URL в оперативной памяти.
type MemoryURLStore struct {
	mx   sync.RWMutex
	urls map[string]UserURL
	cfg  *config.Config
}

// NewMemoryURLStore создаёт новый экземпляр MemoryURLStore.
func NewMemoryURLStore(cfg *config.Config) *MemoryURLStore {
	store := &MemoryURLStore{urls: make(map[string]UserURL), cfg: cfg}
	log := logger.NewLogger()
	err := store.LoadSnapshot()
	if err != nil {
		log.Log.Error("Error loading snapshot", zap.Error(err))
	}
	return store
}

// SetURL сохраняет новый URL в хранилище.
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

// SetBatchURL сохраняет пакет URL в хранилище.
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

// GetURL возвращает оригинальный URL по короткому ключу.
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
	if value.IsDeleted {
		return "", ErrAlreadyDeleted
	}
	return value.URL, nil
}

// GetUserURLs возвращает список URL пользователя.
func (s *MemoryURLStore) GetUserURLs(userID string) ([]models.UserShortURLReadDTO, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if s.urls == nil {
		return nil, ErrNotInitialized
	}
	var readDTO []models.UserShortURLReadDTO
	for key, value := range s.urls {
		if value.UserID == userID && !value.IsDeleted {
			readDTO = append(readDTO, models.UserShortURLReadDTO{ShortURL: fmt.Sprintf("%s/%s", s.cfg.BaseURL, key), OriginalURL: value.URL})
		}
	}
	if len(readDTO) == 0 {
		return nil, ErrNotFound
	}
	return readDTO, nil
}

// DeleteURLs помечает список URL как удалённые.
func (s *MemoryURLStore) DeleteURLs(userID string, urls []string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.urls == nil {
		return ErrNotInitialized
	}
	if len(userID) == 0 {
		return ErrEmptyUserID
	}
	if len(urls) == 0 {
		return ErrEmptyURLs
	}

	for _, shortURL := range urls {
		if value, exists := s.urls[shortURL]; exists {
			if value.UserID == userID {
				value.IsDeleted = true
				s.urls[shortURL] = value
			}
		}
	}

	return nil
}

// Close завершает работу хранилища, создавая снапшот.
func (s *MemoryURLStore) Close() error {
	return s.CreateSnapshot()
}

// CountURLs возвращает количество всех сокращённых URL в памяти.
func (s *MemoryURLStore) CountURLs() (int, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return len(s.urls), nil
}

// CountUsers возвращает количество уникальных пользователей в памяти.
func (s *MemoryURLStore) CountUsers() (int, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	users := make(map[string]struct{})
	for _, v := range s.urls {
		users[v.UserID] = struct{}{}
	}
	return len(users), nil
}
