package mocks

import (
	"fmt"

	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/shekshuev/shortener/internal/app/store"
	"github.com/stretchr/testify/mock"
)

// MockStore - моковая реализация хранилища URL.
type MockStore struct {
	mock.Mock
	urls map[string]store.UserURL
}

// ErrNotFound - ошибка, возникающая при отсутствии запрашиваемого URL.
var ErrNotFound = fmt.Errorf("not found")

// NewURLStore создаёт новый моковый стор для хранения URL.
func NewURLStore() *MockStore {
	store := &MockStore{urls: make(map[string]store.UserURL)}
	return store
}

// SetURL сохраняет URL в хранилище.
func (m *MockStore) SetURL(key, value, userID string) (string, error) {
	if len(key) == 0 {
		return "", store.ErrEmptyKey
	}
	if len(value) == 0 {
		return "", store.ErrEmptyValue
	}
	if len(userID) == 0 {
		return "", store.ErrEmptyUserID
	}
	m.urls[key] = store.UserURL{UserID: userID, URL: value}
	return value, nil
}

// SetBatchURL сохраняет пакет URL в хранилище.
func (m *MockStore) SetBatchURL(createDTO []models.BatchShortURLCreateDTO, userID string) error {
	if len(userID) == 0 {
		return store.ErrEmptyUserID
	}
	if createDTO == nil {
		return store.ErrEmptyValue
	}
	for _, dto := range createDTO {
		if len(dto.ShortURL) == 0 {
			return store.ErrEmptyKey
		}
		if len(dto.OriginalURL) == 0 {
			return store.ErrEmptyValue
		}
		m.urls[dto.ShortURL] = store.UserURL{UserID: userID, URL: dto.OriginalURL}
	}
	return nil
}

// GetURL возвращает оригинальный URL по короткому ключу.
func (m *MockStore) GetURL(key string) (string, error) {
	value, exists := m.urls[key]
	if !exists {
		return "", ErrNotFound
	}
	return value.URL, nil
}

// GetUserURLs возвращает все URL, принадлежащие пользователю.
func (m *MockStore) GetUserURLs(userID string) ([]models.UserShortURLReadDTO, error) {
	var readDTO []models.UserShortURLReadDTO
	for key, value := range m.urls {
		if value.UserID == userID {
			readDTO = append(readDTO, models.UserShortURLReadDTO{ShortURL: key, OriginalURL: value.URL})
		}
	}
	if len(readDTO) == 0 {
		return nil, ErrNotFound
	}
	return readDTO, nil
}

// DeleteURLs помечает список URL как удалённые.
func (m *MockStore) DeleteURLs(userID string, urls []string) error {
	if m.urls == nil {
		return store.ErrNotInitialized
	}
	if len(userID) == 0 {
		return store.ErrEmptyUserID
	}
	if len(urls) == 0 {
		return store.ErrEmptyURLs
	}

	for _, shortURL := range urls {
		if value, exists := m.urls[shortURL]; exists {
			if value.UserID == userID {
				value.IsDeleted = true
				m.urls[shortURL] = value
			}
		}
	}

	return nil
}

// CheckDBConnection проверяет подключение к базе данных (мокается для тестов).
func (m *MockStore) CheckDBConnection() error {
	args := m.Called()
	return args.Error(0)
}

// Close закрывает подключение к хранилищу (мокается для тестов).
func (m *MockStore) Close() error {
	args := m.Called()
	return args.Error(0)
}

// CountURLs возвращает количество всех сокращённых URL в моке.
func (m *MockStore) CountURLs() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// CountUsers возвращает количество уникальных пользователей в моке.
func (m *MockStore) CountUsers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}
