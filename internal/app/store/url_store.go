package store

import "errors"

type URLStore struct {
	urls map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{urls: make(map[string]string)}
}

func (s *URLStore) SetURL(key, value string) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}
	if len(value) == 0 {
		return errors.New("value cannot be empty")
	}
	s.urls[key] = value
	return nil
}

func (s *URLStore) GetURL(key string) (string, error) {
	value, exists := s.urls[key]
	if !exists {
		return "", errors.New("not found")
	}
	return value, nil
}
