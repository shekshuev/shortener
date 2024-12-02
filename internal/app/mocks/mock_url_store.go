package mocks

import (
	"fmt"

	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/shekshuev/shortener/internal/app/store"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
	urls map[string]store.UserURL
}

var ErrNotFound = fmt.Errorf("not found")

func NewURLStore() *MockStore {
	store := &MockStore{urls: make(map[string]store.UserURL)}
	return store
}

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

func (m *MockStore) GetURL(key string) (string, error) {
	value, exists := m.urls[key]
	if !exists {
		return "", ErrNotFound
	}
	return value.URL, nil
}

func (m *MockStore) CheckDBConnection() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStore) Close() error {
	args := m.Called()
	return args.Error(0)
}
