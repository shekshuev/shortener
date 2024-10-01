package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_create(t *testing.T) {
	testCases := []struct {
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, expectedCode: http.StatusCreated, body: "https://ya.ru", expectedBody: "result here"},
		{method: http.MethodPut, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodGet, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodDelete, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodPatch, expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			create(w, r)
			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tc.expectedBody != "" {
				assert.Len(t, w.Body.String(), len(fmt.Sprintf("http://%s/12345678", r.Host)),
					"Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func Test_get(t *testing.T) {
	url := "https://ya.ru"
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(url))
	w := httptest.NewRecorder()
	create(w, r)
	shortedId := "/" + strings.Split(w.Body.String(), "/")[3]
	testCases := []struct {
		method       string
		target       string
		expectedCode int
	}{
		{method: http.MethodGet, target: shortedId, expectedCode: http.StatusTemporaryRedirect},
		{method: http.MethodGet, target: "/unknown", expectedCode: http.StatusBadRequest},
		{method: http.MethodPut, target: shortedId, expectedCode: http.StatusBadRequest},
		{method: http.MethodPost, target: shortedId, expectedCode: http.StatusBadRequest},
		{method: http.MethodDelete, target: shortedId, expectedCode: http.StatusBadRequest},
		{method: http.MethodPatch, target: shortedId, expectedCode: http.StatusBadRequest},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.target, nil)
			w := httptest.NewRecorder()
			get(w, r)
			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if w.Code == http.StatusTemporaryRedirect {
				assert.Equal(t, url, w.Header().Get("Location"), "Адрес не совпадает с ожидаемым")
			}
		})
	}
}
