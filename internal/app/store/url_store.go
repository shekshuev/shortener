package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/logger"
)

type URLStore struct {
	urls map[string]string
	cfg  *config.Config
}

type serializeData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var ErrEmptyKey = fmt.Errorf("key cannot be empty")
var ErrEmptyValue = fmt.Errorf("value cannot be empty")
var ErrNotFound = fmt.Errorf("not found")

func NewURLStore(cfg *config.Config) *URLStore {
	store := &URLStore{urls: make(map[string]string), cfg: cfg}
	if _, err := os.Stat(cfg.FileStoragePath); err == nil {
		file, err := os.Open(cfg.FileStoragePath)
		if err != nil {
			logger.Log.Error("error opening file", zap.Error(err))
			return store
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var urlData serializeData
			err := json.Unmarshal(scanner.Bytes(), &urlData)
			if err != nil {
				logger.Log.Error("error parsing JSON", zap.Error(err))
				continue
			}
			store.urls[urlData.ShortURL] = urlData.OriginalURL
		}
		if err := scanner.Err(); err != nil {
			logger.Log.Error("error scanning file", zap.Error(err))
		}
	}
	return store
}

func (s *URLStore) SetURL(key, value string) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if len(value) == 0 {
		return ErrEmptyValue
	}
	s.urls[key] = value
	return s.persist(key)
}

func (s *URLStore) GetURL(key string) (string, error) {
	value, exists := s.urls[key]
	if !exists {
		return "", ErrNotFound
	}
	return value, nil
}

func (s *URLStore) persist(key string) error {
	file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	urlData := serializeData{
		UUID:        uuid.New().String(),
		ShortURL:    key,
		OriginalURL: s.urls[key],
	}

	data, err := json.Marshal(urlData)
	if err != nil {
		return err
	}

	_, err = file.Write(append(data, '\n'))
	if err != nil {
		return err
	}

	return nil
}
