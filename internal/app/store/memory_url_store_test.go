package store

import (
	"testing"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestMemoryURLStore_SetURL(t *testing.T) {
	testCases := []struct {
		key      string
		value    string
		name     string
		userID   string
		hasError bool
	}{
		{name: "Normal key and value", key: "test", value: "test", userID: "1", hasError: false},
		{name: "Empty key", key: "", value: "test", userID: "1", hasError: true},
		{name: "Empty value", key: "test", value: "", userID: "1", hasError: true},
		{name: "Empty userID", key: "test", value: "test", userID: "", hasError: true},
	}
	cfg := config.GetConfig()
	s := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.SetURL(tc.key, tc.value, tc.userID)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
			}
		})
	}
}

func TestMemoryURLStore_SetBatchURL(t *testing.T) {
	testCases := []struct {
		name      string
		createDTO []models.BatchShortURLCreateDTO
		userID    string
		hasError  bool
	}{
		{name: "Not empty list with correct values", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: "test1"},
			{CorrelationID: "test2", OriginalURL: "https://google.com", ShortURL: "test2"},
		}, userID: "1", hasError: false},
		{name: "Empty list", createDTO: []models.BatchShortURLCreateDTO{}, userID: "1", hasError: false},
		{name: "Not empty list with empty short url", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: ""},
			{CorrelationID: "test2", OriginalURL: "https://google.com", ShortURL: "test2"},
		}, userID: "1", hasError: true},
		{name: "Not empty list with empty original url", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: "test1"},
			{CorrelationID: "test2", OriginalURL: "", ShortURL: "test2"},
		}, userID: "1", hasError: true},
		{name: "Nil list", createDTO: nil, userID: "1", hasError: true},
		{name: "Empty user id", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: "test1"},
			{CorrelationID: "test2", OriginalURL: "", ShortURL: "test2"},
		}, userID: "", hasError: true},
	}
	cfg := config.GetConfig()
	s := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.SetBatchURL(tc.createDTO, tc.userID)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
			}
		})
	}
}

func TestMemoryURLStore_GetURL(t *testing.T) {
	testCases := []struct {
		key    string
		getKey string
		value  string
		name   string
		userID string
	}{
		{name: "Get existing value", key: "test", getKey: "test", value: "test", userID: "1"},
		{name: "Get not existing value", key: "test", getKey: "not exists", value: "test", userID: "1"},
	}
	cfg := config.GetConfig()
	s := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.SetURL(tc.key, tc.value, tc.userID)
			assert.Nil(t, err, "Set error is not nil")
			res, err := s.GetURL(tc.getKey)
			if tc.key == tc.getKey {
				assert.Equal(t, res, tc.value, "Get result is not equal to test value")
				assert.Nil(t, err, "Get error is not nil")
			} else {
				assert.Len(t, res, 0, "Get result is not nil")
				assert.NotNil(t, err, "Get error is nil")
			}
		})
	}
}
