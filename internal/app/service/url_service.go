package service

import (
	"errors"
	"fmt"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/shekshuev/shortener/internal/app/store"
	"github.com/shekshuev/shortener/internal/utils"
)

// Service - интерфейс для работы с URL.
type Service interface {
	CreateShortURL(longURL, userID string) (string, error)
	BatchCreateShortURL(createDTO []models.BatchShortURLCreateDTO, userID string) ([]models.BatchShortURLReadDTO, error)
	GetLongURL(shortURL string) (string, error)
	GetUserURLs(userID string) ([]models.UserShortURLReadDTO, error)
	DeleteURLs(userID string, urls []string)
	CheckDBConnection() error
}

// URLService - реализация сервиса для управления URL.
type URLService struct {
	store store.URLStore
	cfg   *config.Config
}

// ErrNotPostgresStore - ошибка, указывающая на использование in-memory хранилища вместо Postgres.
var ErrNotPostgresStore = fmt.Errorf("app using in-memory store, not postgres")

// NewURLService создаёт новый экземпляр URLService.
func NewURLService(store store.URLStore, cfg *config.Config) *URLService {
	return &URLService{store: store, cfg: cfg}
}

// ErrFailedToShorten - ошибка при создании короткого URL.
var ErrFailedToShorten = fmt.Errorf("failed to create short url")

// CreateShortURL создаёт короткий URL.
func (s *URLService) CreateShortURL(longURL, userID string) (string, error) {
	shorted, err := utils.Shorten(longURL)
	if err != nil {
		return "", ErrFailedToShorten
	}
	shortURL, err := s.store.SetURL(shorted, longURL, userID)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortURL), err
		}
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shorted), nil
}

// BatchCreateShortURL создаёт несколько коротких URL в пакете.
func (s *URLService) BatchCreateShortURL(createDTO []models.BatchShortURLCreateDTO, userID string) ([]models.BatchShortURLReadDTO, error) {
	for i := 0; i < len(createDTO); i++ {
		shorted, err := utils.Shorten(createDTO[i].OriginalURL)
		if err != nil {
			return nil, ErrFailedToShorten
		}
		createDTO[i].ShortURL = shorted
	}

	err := s.store.SetBatchURL(createDTO, userID)
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

// GetLongURL возвращает оригинальный URL по короткому.
func (s *URLService) GetLongURL(shortURL string) (string, error) {
	longURL, err := s.store.GetURL(shortURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

// GetUserURLs возвращает список URL пользователя.
func (s *URLService) GetUserURLs(userID string) ([]models.UserShortURLReadDTO, error) {
	readDTO, err := s.store.GetUserURLs(userID)
	if err != nil {
		return nil, err
	}
	return readDTO, nil
}

// DeleteURLs удаляет список URL пользователя.
func (s *URLService) DeleteURLs(userID string, urls []string) {
	s.store.DeleteURLs(userID, urls)
}

// CheckDBConnection проверяет соединение с базой данных.
func (s *URLService) CheckDBConnection() error {
	if dbChecker, ok := s.store.(store.DatabaseChecker); ok {
		return dbChecker.CheckDBConnection()
	}
	return ErrNotPostgresStore
}
