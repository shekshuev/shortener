package store

import "fmt"

type URLStore interface {
	SetURL(key, value string) error
	GetURL(key string) (string, error)
	Close() error
}

type DatabaseChecker interface {
	CheckDBConnection() error
}

var ErrEmptyKey = fmt.Errorf("key cannot be empty")
var ErrEmptyValue = fmt.Errorf("value cannot be empty")
var ErrNotFound = fmt.Errorf("not found")
var ErrNotInitialized = fmt.Errorf("store not initialized")
