package store

type URLStore interface {
	SetURL(key, value string) error
	GetURL(key string) (string, error)
	CheckDBConnection() error
}
