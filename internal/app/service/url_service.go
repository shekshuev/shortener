package service

import (
	"errors"
	"fmt"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/shekshuev/shortener/internal/app/store"
	"github.com/shekshuev/shortener/internal/utils"
)

type Service interface {
	CreateShortURL(longURL string) (string, error)
	BatchCreateShortURL(createDTO []models.BatchShortURLCreateDTO) ([]models.BatchShortURLReadDTO, error)
	GetLongURL(shortURL string) (string, error)
	CheckDBConnection() error
}

type URLService struct {
	store store.URLStore
	cfg   *config.Config
}

var ErrNotPostgresStore = fmt.Errorf("app using in-memory store, not postgres")

func NewURLService(store store.URLStore, cfg *config.Config) *URLService {
	return &URLService{store: store, cfg: cfg}
}

var ErrFailedToShorten = fmt.Errorf("failed to create short url")

func (s *URLService) CreateShortURL(longURL string) (string, error) {
	shorted, err := utils.Shorten(longURL)
	if err != nil {
		return "", ErrFailedToShorten
	}
	shortURL, err := s.store.SetURL(shorted, longURL)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortURL), err
		}
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shorted), nil
}

func (s *URLService) BatchCreateShortURL(createDTO []models.BatchShortURLCreateDTO) ([]models.BatchShortURLReadDTO, error) {
	for i := 0; i < len(createDTO); i++ {
		shorted, err := utils.Shorten(createDTO[i].OriginalURL)
		if err != nil {
			return nil, ErrFailedToShorten
		}
		createDTO[i].ShortURL = shorted
	}

	err := s.store.SetBatchURL(createDTO)
	readDTO := make([]models.BatchShortURLReadDTO, 0, len(createDTO))
	for _, dto := range createDTO {
		readDTO = append(readDTO, models.BatchShortURLReadDTO{
			CorrelationID: dto.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", s.cfg.BaseURL, dto.ShortURL),
		})
	}
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return readDTO, err
		}
		return nil, err
	}

	return readDTO, nil
}

func (s *URLService) GetLongURL(shortURL string) (string, error) {
	longURL, err := s.store.GetURL(shortURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (s *URLService) CheckDBConnection() error {
	if dbChecker, ok := s.store.(store.DatabaseChecker); ok {
		return dbChecker.CheckDBConnection()
	}
	return ErrNotPostgresStore
}
