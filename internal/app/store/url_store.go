package store

import "errors"

type URLStore struct {
	urls map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{urls: make(map[string]string)}
}

func (s *URLStore) SetURL(key, value string) {
	s.urls[key] = value
}

func (s *URLStore) GetURL(key string) (string, error) {
	value, exists := s.urls[key]
	if !exists {
		return "", errors.New("not found")
	}
	return value, nil
}
