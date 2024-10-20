package utils

import (
	"fmt"
	"math/rand"
)

const (
	ShortenLength = 8
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var ErrEmptyString = fmt.Errorf("string should not be empty")

func Shorten(s string) (string, error) {
	if len(s) == 0 {
		return "", ErrEmptyString
	}
	keys := make([]rune, ShortenLength)
	for i := range ShortenLength {
		keys[i] = rune(charset[rand.Intn(len(charset))])
	}
	return string(keys), nil
}
