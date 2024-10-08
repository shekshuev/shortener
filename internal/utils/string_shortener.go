package utils

import (
	"errors"
	"math/rand"
)

const (
	ShortenLength = 8
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func Shorten(s string) (string, error) {
	if len(s) == 0 {
		return "", errors.New("string should not be empty")
	}
	keys := make([]rune, ShortenLength)
	for i := range ShortenLength {
		keys[i] = rune(charset[rand.Intn(len(charset))])
	}
	return string(keys), nil
}
