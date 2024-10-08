package service

import (
	"errors"

	"github.com/shekshuev/shortener/internal/app/store"
	"github.com/shekshuev/shortener/internal/utils"
)

type URLService struct {
	store *store.URLStore
}

func NewURLService(store *store.URLStore) *URLService {
	return &URLService{store: store}
}

func (s *URLService) CreateShortURL(longURL string) (string, error) {
	shortURL, err := utils.Shorten(longURL)
	if err != nil {
		return "", errors.New("failed to create short url")
	}
	s.store.SetURL(shortURL, longURL)
	return shortURL, nil
}

func (s *URLService) GetLongURL(shortURL string) (string, error) {
	longURL, err := s.store.GetURL(shortURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}
