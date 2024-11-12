package mocks

import (
	"fmt"

	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
	urls map[string]string
}

var ErrNotFound = fmt.Errorf("not found")

func NewURLStore() *MockStore {
	store := &MockStore{urls: make(map[string]string)}
	return store
}

func (m *MockStore) SetURL(key, value string) error {
	m.urls[key] = value
	return nil
}

func (m *MockStore) SetBatchURL(createDTO []models.BatchShortURLCreateDTO) error {
	for _, dto := range createDTO {
		m.urls[dto.ShortURL] = dto.OriginalURL
	}
	return nil
}

func (m *MockStore) GetURL(key string) (string, error) {
	value, exists := m.urls[key]
	if !exists {
		return "", ErrNotFound
	}
	return value, nil
}

func (m *MockStore) CheckDBConnection() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStore) Close() error {
	args := m.Called()
	return args.Error(0)
}
