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
		longURL string
		name    string
	}{
		{name: "Normal long URL", longURL: "https://example.com"},
		{name: "Empty long URL", longURL: ""},
	}
	cfg := config.GetConfig()
	s := mocks.NewURLStore()
	service := NewURLService(s, &cfg)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortURL, err := service.CreateShortURL(tc.longURL)
			if len(tc.longURL) == 0 {
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
		hasError  bool
	}{
		{name: "Not empty list with correct values", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru"},
			{CorrelationID: "test2", OriginalURL: "https://google.com"},
		}, hasError: false},
		{name: "Empty list", createDTO: []models.BatchShortURLCreateDTO{}, hasError: false},
		{name: "Not empty list with empty original url", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru"},
			{CorrelationID: "test2", OriginalURL: ""},
		}, hasError: true},
		{name: "Nil list", createDTO: nil, hasError: false},
	}
	cfg := config.GetConfig()
	s := mocks.NewURLStore()
	service := NewURLService(s, &cfg)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			readDTO, err := service.BatchCreateShortURL(tc.createDTO)
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
	testCases := []struct {
		shorted string
		name    string
	}{
		{name: "Existing short URL", shorted: shorted},
		{name: "Non-existing short URL", shorted: "non-existing"},
	}
	cfg := config.GetConfig()
	s := mocks.NewURLStore()
	service := NewURLService(s, &cfg)
	err := s.SetURL(shorted, longURL)
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
