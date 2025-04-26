package service

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/mocks"
	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestNewURLService(t *testing.T) {
	t.Run("Test NewURLService", func(t *testing.T) {
		cfg := config.GetConfig()
		s := mocks.NewURLStore()
		service := NewURLService(s, &cfg)
		assert.Equal(t, service.store, s, "URLService has incorrect store")
	})
}

func TestURLService_CreateShortURL(t *testing.T) {
	shorted := "12345678"
	testCases := []struct {
		longURL  string
		name     string
		userID   string
		hasError bool
	}{
		{name: "Normal long URL", longURL: "https://example.com", userID: "1", hasError: false},
		{name: "Empty long URL", longURL: "", userID: "1", hasError: true},
		{name: "Empty userID", longURL: "https://example.com", userID: "", hasError: true},
	}
	cfg := config.GetConfig()
	s := mocks.NewURLStore()
	service := NewURLService(s, &cfg)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortURL, err := service.CreateShortURL(tc.longURL, tc.userID)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Len(t, shortURL, len(fmt.Sprintf("%s/%s", cfg.BaseURL, shorted)))
			}
		})
	}
}

func TestURLService_BatchCreateShortURL(t *testing.T) {
	shorted := "12345678"
	testCases := []struct {
		name      string
		createDTO []models.BatchShortURLCreateDTO
		userID    string
		hasError  bool
	}{
		{name: "Not empty list with correct values", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru"},
			{CorrelationID: "test2", OriginalURL: "https://google.com"},
		}, userID: "1", hasError: false},
		{name: "Empty list", createDTO: []models.BatchShortURLCreateDTO{}, userID: "1", hasError: false},
		{name: "Not empty list with empty original url", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru"},
			{CorrelationID: "test2", OriginalURL: ""},
		}, userID: "1", hasError: true},
		{name: "Nil list", createDTO: nil, userID: "1", hasError: true},
		{name: "Empty user id", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru"},
			{CorrelationID: "test2", OriginalURL: ""},
		}, userID: "", hasError: true},
	}
	cfg := config.GetConfig()
	s := mocks.NewURLStore()
	service := NewURLService(s, &cfg)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			readDTO, err := service.BatchCreateShortURL(tc.createDTO, tc.userID)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")

			}
			for _, dto := range readDTO {
				assert.Len(t, dto.ShortURL, len(fmt.Sprintf("%s/%s", cfg.BaseURL, shorted)))
			}
		})
	}
}

func TestURLService_GetLongURL(t *testing.T) {
	longURL := "https://example.com"
	shorted := "12345678"
	userID := "1"
	testCases := []struct {
		shorted  string
		name     string
		hasError bool
	}{
		{name: "Existing short URL", shorted: shorted, hasError: false},
		{name: "Non-existing short URL", shorted: "non-existing", hasError: true},
	}
	cfg := config.GetConfig()
	s := mocks.NewURLStore()
	service := NewURLService(s, &cfg)
	_, err := s.SetURL(shorted, longURL, userID)
	assert.Nil(t, err, "Set url store error is not nil")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			longURL, err := service.GetLongURL(tc.shorted)
			if tc.shorted == "non-existing" {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Equal(t, longURL, longURL)
			}
		})
	}
}

func TestURLService_GetUserURLs(t *testing.T) {
	longURL := "https://example.com"
	shorted := "12345678"
	userID := "1"
	testCases := []struct {
		shorted   string
		name      string
		getUserID string
		hasError  bool
	}{
		{name: "Get with correct userID", shorted: shorted, getUserID: userID, hasError: false},
		{name: "Get with wrong userID", shorted: "non-existing", getUserID: "2", hasError: true},
	}
	cfg := config.GetConfig()
	s := mocks.NewURLStore()
	service := NewURLService(s, &cfg)
	_, err := s.SetURL(shorted, longURL, userID)
	assert.Nil(t, err, "Set url store error is not nil")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			readDTO, err := service.GetUserURLs(tc.getUserID)
			if tc.shorted == "non-existing" {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
				assert.Len(t, readDTO, 1)
			}
		})
	}
}

func TestURLService_CheckDBConnection(t *testing.T) {
	testCases := []struct {
		name     string
		hasError bool
		error    error
	}{
		{name: "Success", hasError: false, error: nil},
		{name: "Error", hasError: true, error: sql.ErrConnDone},
	}

	cfg := config.GetConfig()

	mockStore := mocks.NewURLStore()
	service := NewURLService(mockStore, &cfg)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStore.On("CheckDBConnection").Return(tc.error)
			err := service.CheckDBConnection()
			assert.Equal(t, tc.hasError, err != nil, "CheckDBConnection failed")
			mockStore.AssertCalled(t, "CheckDBConnection")
			mockStore.ExpectedCalls = nil
		})
	}
}

func TestURLService_GetStats(t *testing.T) {
	cfg := config.GetConfig()
	mockStore := mocks.NewURLStore()
	service := NewURLService(mockStore, &cfg)

	testCases := []struct {
		name          string
		urlsCount     int
		usersCount    int
		urlsError     error
		usersError    error
		expectError   bool
		expectedURLs  int
		expectedUsers int
	}{
		{
			name:          "Success getting stats",
			urlsCount:     10,
			usersCount:    5,
			urlsError:     nil,
			usersError:    nil,
			expectError:   false,
			expectedURLs:  10,
			expectedUsers: 5,
		},
		{
			name:          "Error getting URLs count",
			urlsCount:     0,
			usersCount:    0,
			urlsError:     fmt.Errorf("db error"),
			usersError:    nil,
			expectError:   true,
			expectedURLs:  0,
			expectedUsers: 0,
		},
		{
			name:          "Error getting Users count",
			urlsCount:     10,
			usersCount:    0,
			urlsError:     nil,
			usersError:    fmt.Errorf("db error"),
			expectError:   true,
			expectedURLs:  0,
			expectedUsers: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStore.On("CountURLs").Return(tc.urlsCount, tc.urlsError)
			mockStore.On("CountUsers").Return(tc.usersCount, tc.usersError)

			stats, err := service.GetStats()
			if tc.expectError {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Equal(t, 0, stats.URLs)
				assert.Equal(t, 0, stats.Users)
			} else {
				assert.Nil(t, err, "Unexpected error")
				assert.Equal(t, tc.expectedURLs, stats.URLs, "URLs count mismatch")
				assert.Equal(t, tc.expectedUsers, stats.Users, "Users count mismatch")
			}

			mockStore.AssertCalled(t, "CountURLs")
			if tc.urlsError == nil {
				mockStore.AssertCalled(t, "CountUsers")
			}
			mockStore.ExpectedCalls = nil
		})
	}
}
