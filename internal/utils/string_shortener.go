package utils

import (
	"fmt"
	"math/rand"
)

const (
	// ShortenLength Длина сокращённого URL
	ShortenLength = 8
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// ErrEmptyString - ошибка, возникающая при попытке сократить пустую строку.
var ErrEmptyString = fmt.Errorf("string should not be empty")

// Shorten генерирует случайную строку длиной ShortenLength для сокращения URL.
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
