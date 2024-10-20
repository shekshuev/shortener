package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_shorten(t *testing.T) {
	expectedLength := 8
	testCases := []struct {
		name           string
		input          string
		expectedLength int
	}{
		{name: "Normal URL", input: "https://ya.ru", expectedLength: expectedLength},
		{name: "Long string", input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua", expectedLength: expectedLength},
		{name: "Short string", input: "_", expectedLength: expectedLength},
		{name: "Empty string, error is not nil", input: "", expectedLength: 0},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Shorten(tc.input)
			assert.Equal(t, len(res), tc.expectedLength, "Wrong length")
			if len(res) == 0 {
				assert.NotNil(t, err, "Error is nil")
			}
		})
	}
}
