package service

import (
	"errors"
	"fmt"
	"github.com/shekshuev/shortener/internal/app/config"

	"github.com/shekshuev/shortener/internal/app/store"
	"github.com/shekshuev/shortener/internal/utils"
)

type URLService struct {
	store *store.URLStore
	cfg   *config.Config
}

func NewURLService(store *store.URLStore, cfg *config.Config) *URLService {
	return &URLService{store: store, cfg: cfg}
}

func (s *URLService) CreateShortURL(longURL string) (string, error) {
	shorted, err := utils.Shorten(longURL)
	if err != nil {
		return "", errors.New("failed to create short url")
	}
	err = s.store.SetURL(shorted, longURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shorted), nil
}

func (s *URLService) GetLongURL(shortURL string) (string, error) {
	longURL, err := s.store.GetURL(shortURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}
