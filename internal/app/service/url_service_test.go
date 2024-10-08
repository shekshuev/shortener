package service

import (
	"fmt"
	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewURLService(t *testing.T) {
	t.Run("Test NewURLService", func(t *testing.T) {
		s := store.NewURLStore()
		cfg := config.GetConfig()
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
	s := store.NewURLStore()
	cfg := config.GetConfig()
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
	s := store.NewURLStore()
	cfg := config.GetConfig()
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
