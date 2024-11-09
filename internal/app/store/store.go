package store

type Store interface {
	SetURL(key, value string) error
	GetURL(key string) (string, error)
	CheckDBConnection() error
}
