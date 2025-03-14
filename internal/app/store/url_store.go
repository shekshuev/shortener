package store

import (
	"fmt"

	"github.com/shekshuev/shortener/internal/app/models"
)

// URLStore - интерфейс для работы с хранилищем URL.
type URLStore interface {
	SetURL(key, value, userID string) (string, error)
	SetBatchURL(createDTO []models.BatchShortURLCreateDTO, userID string) error
	GetURL(key string) (string, error)
	GetUserURLs(userID string) ([]models.UserShortURLReadDTO, error)
	DeleteURLs(userID string, urls []string) error
	Close() error
}

// DatabaseChecker - интерфейс для проверки соединения с базой данных.
type DatabaseChecker interface {
	CheckDBConnection() error
}

// Общие ошибки, возникающие при работе с хранилищем.
var (
	ErrAlreadyExists  = fmt.Errorf("url already exists")      // Ошибка: URL уже существует
	ErrEmptyKey       = fmt.Errorf("key cannot be empty")     // Ошибка: ключ не может быть пустым
	ErrEmptyValue     = fmt.Errorf("value cannot be empty")   // Ошибка: значение не может быть пустым
	ErrEmptyUserID    = fmt.Errorf("user ID cannot be empty") // Ошибка: идентификатор пользователя не может быть пустым
	ErrNotFound       = fmt.Errorf("not found")               // Ошибка: запись не найдена
	ErrNotInitialized = fmt.Errorf("store not initialized")   // Ошибка: хранилище не инициализировано
	ErrEmptyURLs      = fmt.Errorf("no urls provided")        // Ошибка: список URL пуст
	ErrAlreadyDeleted = fmt.Errorf("urls already deleted")    // Ошибка: URL уже удалены
)
