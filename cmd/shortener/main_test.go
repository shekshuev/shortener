package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/shekshuev/shortener/internal/app/config"
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
		{method: http.MethodPut, expectedCode: http.StatusBadRequest},
		{method: http.MethodGet, expectedCode: http.StatusBadRequest},
		{method: http.MethodDelete, expectedCode: http.StatusBadRequest},
		{method: http.MethodPatch, expectedCode: http.StatusBadRequest},
	}
	srv := httptest.NewServer(http.HandlerFunc(create))

	config.SetConfig(srv.URL, srv.URL)

	defer srv.Close()

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL
			resp, err := req.SetBody(tc.body).Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tc.expectedBody != "" {
				assert.Len(t, string(resp.Body()), len(srv.URL)+9, "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func Test_get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", create)
	mux.HandleFunc("/{shorted}", get)
	srv := httptest.NewServer(mux)

	config.SetConfig(srv.URL, srv.URL)

	defer srv.Close()

	url := "https://ya.ru"

	resp, err := resty.New().R().SetBody(url).Post(srv.URL)
	assert.NoError(t, err, "error making HTTP request")
	shortedID := "/" + strings.Split(string(resp.Body()), "/")[3]

	testCases := []struct {
		method              string
		target              string
		expectedCode        int
		expectedRedirectURL string
	}{
		{method: http.MethodGet, target: shortedID, expectedCode: http.StatusOK, expectedRedirectURL: url},
		{method: http.MethodGet, target: "/unknown", expectedCode: http.StatusBadRequest},
		{method: http.MethodPut, target: shortedID, expectedCode: http.StatusBadRequest},
		{method: http.MethodPost, target: shortedID, expectedCode: http.StatusBadRequest},
		{method: http.MethodDelete, target: shortedID, expectedCode: http.StatusBadRequest},
		{method: http.MethodPatch, target: shortedID, expectedCode: http.StatusBadRequest},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.target

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			if len(tc.expectedRedirectURL) > 0 {
				assert.Contains(t, resp.RawResponse.Request.URL.String(), url, "Wrong redirect url")
			}
		})
	}
}
