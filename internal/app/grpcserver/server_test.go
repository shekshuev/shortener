package grpcserver

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/mocks"
	"github.com/shekshuev/shortener/internal/app/proto"
	"github.com/shekshuev/shortener/internal/app/service"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *Server {
	cfg := config.GetConfig()
	s := mocks.NewURLStore()
	srv := NewServer(service.NewURLService(s, &cfg))
	return srv
}

func TestServer_Shorten(t *testing.T) {
	srv := setupTestServer()
	ctx := context.Background()

	testCases := []struct {
		name         string
		request      *proto.ShortenRequest
		expectError  bool
		expectResult bool
	}{
		{name: "Success", request: &proto.ShortenRequest{Url: "https://example.com", UserId: "test-user-id"}, expectError: false, expectResult: true},
		{name: "Empty URL", request: &proto.ShortenRequest{Url: "", UserId: "test-user-id"}, expectError: true, expectResult: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := srv.Shorten(ctx, tc.request)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Contains(t, resp.Result, "http")
			}
		})
	}
}

func TestServer_BatchShorten(t *testing.T) {
	srv := setupTestServer()
	ctx := context.Background()

	testCases := []struct {
		name        string
		request     *proto.BatchShortenRequest
		expectError bool
		expectedLen int
	}{
		{
			name: "Success",
			request: &proto.BatchShortenRequest{
				Items: []*proto.BatchShortenRequestItem{
					{CorrelationId: "id1", OriginalUrl: "https://example.com"},
					{CorrelationId: "id2", OriginalUrl: "https://golang.org"},
				},
				UserId: "test-user-id",
			},
			expectError: false, expectedLen: 2,
		},
		{
			name:        "Empty Items",
			request:     &proto.BatchShortenRequest{Items: []*proto.BatchShortenRequestItem{}, UserId: "test-user-id"},
			expectError: false, expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := srv.BatchShorten(ctx, tc.request)
			assert.NoError(t, err)
			assert.Len(t, resp.Items, tc.expectedLen)
		})
	}
}

func TestServer_GetUserURLs(t *testing.T) {
	srv := setupTestServer()
	ctx := context.Background()

	t.Run("Get User URLs", func(t *testing.T) {
		_, _ = srv.Shorten(ctx, &proto.ShortenRequest{Url: "https://example.com", UserId: "test-user-id"})
		resp, err := srv.GetUserURLs(ctx, &proto.UserURLsRequest{UserId: "test-user-id"})
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Urls)
	})
}

func TestServer_DeleteUserURLs(t *testing.T) {
	srv := setupTestServer()
	ctx := context.Background()

	t.Run("Delete User URLs", func(t *testing.T) {
		resp, err := srv.DeleteUserURLs(ctx, &proto.DeleteURLsRequest{ShortUrls: []string{"id1", "id2"}})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestServer_Ping(t *testing.T) {
	cfg := config.GetConfig()
	mockStore := new(mocks.MockStore)
	srv := NewServer(service.NewURLService(mockStore, &cfg))
	ctx := context.Background()

	testCases := []struct {
		name        string
		mockError   error
		expectError bool
	}{
		{name: "Success", mockError: nil, expectError: false},
		{name: "Failure", mockError: sql.ErrConnDone, expectError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStore.ExpectedCalls = nil
			mockStore.On("CheckDBConnection").Return(tc.mockError)
			resp, err := srv.Ping(ctx, &proto.PingRequest{})
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestServer_GetStats(t *testing.T) {
	cfg := config.GetConfig()
	mockStore := new(mocks.MockStore)
	srv := NewServer(service.NewURLService(mockStore, &cfg))
	ctx := context.Background()

	testCases := []struct {
		name            string
		countURLs       int
		countUsers      int
		countURLError   error
		countUsersError error
		expectError     bool
	}{
		{name: "Success", countURLs: 10, countUsers: 5, expectError: false},
		{name: "CountURLs error", countURLError: sql.ErrConnDone, expectError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStore.ExpectedCalls = nil
			if tc.countURLError != nil {
				mockStore.On("CountURLs").Return(0, tc.countURLError)
			} else {
				mockStore.On("CountURLs").Return(tc.countURLs, nil)
				mockStore.On("CountUsers").Return(tc.countUsers, nil)
			}

			resp, err := srv.GetStats(ctx, &proto.StatsRequest{})
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int32(tc.countURLs), resp.Urls)
				assert.Equal(t, int32(tc.countUsers), resp.Users)
			}
		})
	}
}

func TestServer_GetOriginalURL(t *testing.T) {
	srv := setupTestServer()
	ctx := context.Background()

	shortResp, _ := srv.Shorten(ctx, &proto.ShortenRequest{Url: "https://example.com", UserId: "test-user-id"})
	shortID := strings.TrimPrefix(shortResp.Result, config.GetConfig().BaseURL+"/")

	testCases := []struct {
		name        string
		shortURL    string
		expectedURL string
		expectError bool
	}{
		{
			name:        "Success",
			shortURL:    shortID,
			expectedURL: "https://example.com",
			expectError: false,
		},
		{
			name:        "Failure",
			shortURL:    "nonexistent",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := srv.GetOriginalURL(ctx, &proto.GetOriginalURLRequest{ShortUrl: tc.shortURL})

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedURL, resp.OriginalUrl)
			}
		})
	}
}
