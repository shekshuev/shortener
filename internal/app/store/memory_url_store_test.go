package store

import (
	"testing"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func TestMemoryURLStore_SetURL(t *testing.T) {
	testCases := []struct {
		key   string
		value string
		name  string
	}{
		{name: "Normal key and value", key: "test", value: "test"},
		{name: "Empty key", key: "", value: "test"},
		{name: "Empty value", key: "test", value: ""},
	}
	cfg := config.GetConfig()
	s := &MemoryURLStore{urls: make(map[string]string), cfg: &cfg}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.SetURL(tc.key, tc.value)
			if len(tc.key) == 0 || len(tc.value) == 0 {
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
	}{
		{name: "Get existing value", key: "test", getKey: "test", value: "test"},
		{name: "Get not existing value", key: "test", getKey: "not exists", value: "test"},
	}
	cfg := config.GetConfig()
	s := &MemoryURLStore{urls: make(map[string]string), cfg: &cfg}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.SetURL(tc.key, tc.value)
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
