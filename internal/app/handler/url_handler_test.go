package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/service"
	"github.com/shekshuev/shortener/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func TestNewURLHandler(t *testing.T) {
	t.Run("Test NewURLHandler", func(t *testing.T) {
		s := store.NewURLStore()
		cfg := config.GetConfig()
		srv := service.NewURLService(s, &cfg)
		handler := NewURLHandler(srv)
		assert.Equal(t, handler.service, srv, "URLHandler has incorrect service")
	})
}

func TestURLHandler_createURLHandler(t *testing.T) {
	shortedLenWithSlash := 9
	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
	}{
		{name: "Correct body", method: http.MethodPost, expectedCode: http.StatusCreated, body: "https://ya.ru"},
		{name: "Empty body", method: http.MethodPost, expectedCode: http.StatusBadRequest, body: ""},
	}
	s := store.NewURLStore()
	cfg := config.GetConfig()
	srv := service.NewURLService(s, &cfg)
	handler := NewURLHandler(srv)
	httpSrv := httptest.NewServer(handler.Router)

	defer httpSrv.Close()

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL
			resp, err := req.SetBody(tc.body).Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			if len(tc.body) > 0 {
				assert.Len(t, string(resp.Body()), len(cfg.BaseURL)+shortedLenWithSlash, "Wrong body")
			}
		})
	}
}

func TestURLHandler_getURLHandler(t *testing.T) {
	s := store.NewURLStore()
	cfg := config.GetConfig()
	srv := service.NewURLService(s, &cfg)
	handler := NewURLHandler(srv)
	httpSrv := httptest.NewServer(handler.Router)

	defer httpSrv.Close()

	testURL := "https://ya.ru"

	resp, err := resty.New().R().SetBody(testURL).Post(httpSrv.URL)
	assert.NoError(t, err, "error making HTTP request")
	shortedID := strings.Split(string(resp.Body()), "/")[3]

	testCases := []struct {
		method              string
		target              string
		expectedCode        int
		expectedRedirectURL string
	}{
		{method: http.MethodGet, target: "/" + shortedID, expectedCode: http.StatusOK, expectedRedirectURL: testURL},
		{method: http.MethodGet, target: "/unknown", expectedCode: http.StatusBadRequest},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.target

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			if len(tc.expectedRedirectURL) > 0 {
				assert.NotEmpty(t, resp.RawResponse.Request.URL.String(), "Empty redirect url")
			}
		})
	}
}
