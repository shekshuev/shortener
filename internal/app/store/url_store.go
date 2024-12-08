package store

import (
	"fmt"

	"github.com/shekshuev/shortener/internal/app/models"
)

type URLStore interface {
	SetURL(key, value string) (string, error)
	SetBatchURL(createDTO []models.BatchShortURLCreateDTO) error
	GetURL(key string) (string, error)
	Close() error
}

type DatabaseChecker interface {
	CheckDBConnection() error
}

var ErrAlreadyExists = fmt.Errorf("url already exists")
var ErrEmptyKey = fmt.Errorf("key cannot be empty")
var ErrEmptyValue = fmt.Errorf("value cannot be empty")
var ErrNotFound = fmt.Errorf("not found")
var ErrNotInitialized = fmt.Errorf("store not initialized")
