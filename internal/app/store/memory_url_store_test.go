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
		key       string
		getKey    string
		value     string
		name      string
		userID    string
		getUserID string
	}{
		{name: "Get existing value", key: "test", getKey: "test", value: "test", userID: "1", getUserID: "1"},
		{name: "Get not existing value", key: "test", getKey: "not exists", value: "test", userID: "1", getUserID: "1"},
		{name: "Get existing value with wrong userID", key: "test", getKey: "test", value: "test", userID: "1", getUserID: "2"},
	}
	cfg := config.GetConfig()
	s := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.SetURL(tc.key, tc.value, tc.userID)
			assert.Nil(t, err, "Set error is not nil")
			res, err := s.GetURL(tc.getKey, tc.getUserID)
			if tc.key == tc.getKey && tc.userID == tc.getUserID {
				assert.Equal(t, res, tc.value, "Get result is not equal to test value")
				assert.Nil(t, err, "Get error is not nil")
			} else {
				assert.Len(t, res, 0, "Get result is not nil")
				assert.NotNil(t, err, "Get error is nil")
			}
		})
	}
}

func TestMemoryURLStore_GetUserURLs(t *testing.T) {
	userID := "1"
	testCases := []struct {
		name        string
		getUserID   string
		hasError    bool
		originalURL string
		shortURL    string
	}{
		{name: "Get with correct userID", getUserID: userID, hasError: false, originalURL: "https://ya.ru", shortURL: "test1"},
		{name: "Get with wrong userID", getUserID: "2", hasError: true, originalURL: "https://ya.ru", shortURL: "test1"},
	}
	cfg := config.GetConfig()
	s := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.SetURL(tc.shortURL, tc.originalURL, userID)
			assert.Nil(t, err, "Set error is not nil")
			res, err := s.GetUserURLs(tc.getUserID)
			if !tc.hasError {
				assert.Len(t, res, 1, "Get result length is not equal to test value")
				assert.Nil(t, err, "Get error is not nil")
			} else {
				assert.Len(t, res, 0, "Get result is not nil")
				assert.NotNil(t, err, "Get error is nil")
			}
		})
	}
}

func TestMemoryURLStore_DeleteURLs(t *testing.T) {
	userID := "1"
	testCases := []struct {
		name         string
		setupURLs    []models.BatchShortURLCreateDTO
		deleteURLs   []string
		expectedURLs int
		userID       string
		hasError     bool
	}{
		{
			name: "Delete existing URLs",
			setupURLs: []models.BatchShortURLCreateDTO{
				{OriginalURL: "https://ya.ru", ShortURL: "short1"},
				{OriginalURL: "https://google.com", ShortURL: "short2"},
			},
			deleteURLs:   []string{"short1"},
			expectedURLs: 1,
			userID:       userID,
			hasError:     false,
		},
		{
			name: "Delete non-existent URLs",
			setupURLs: []models.BatchShortURLCreateDTO{
				{OriginalURL: "https://ya.ru", ShortURL: "short1"},
			},
			deleteURLs:   []string{"short3"},
			expectedURLs: 1,
			userID:       userID,
			hasError:     false,
		},
		{
			name: "Delete URLs with wrong userID",
			setupURLs: []models.BatchShortURLCreateDTO{
				{OriginalURL: "https://ya.ru", ShortURL: "short1"},
			},
			deleteURLs:   []string{"short1"},
			expectedURLs: 1,
			userID:       "2",
			hasError:     false,
		},
		{
			name:         "Empty delete list",
			setupURLs:    []models.BatchShortURLCreateDTO{},
			deleteURLs:   []string{},
			expectedURLs: 0,
			userID:       userID,
			hasError:     true,
		},
	}
	cfg := config.GetConfig()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}
			err := s.SetBatchURL(tc.setupURLs, userID)
			assert.Nil(t, err, "Error during setup")

			err = s.DeleteURLs(tc.userID, tc.deleteURLs)
			if tc.hasError {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error during deletion")
			}

			userURLs, err := s.GetUserURLs(userID)
			if tc.hasError {
				assert.NotNil(t, err, "Expected error for empty URLs but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error during GetUserURLs")
			}
			assert.Len(t, userURLs, tc.expectedURLs, "Unexpected number of URLs")
		})
	}
}
