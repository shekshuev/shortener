package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestPostgresURLStore_SetURL(t *testing.T) {
	testCases := []struct {
		key      string
		value    string
		name     string
		userID   string
		hasError bool
	}{
		{name: "Normal key and value", key: "test", value: "test", userID: "1", hasError: false},
		{name: "Empty key", key: "", value: "test", userID: "1", hasError: true},
		{name: "Empty value", key: "test", value: "", userID: "1", hasError: true},
		{name: "Empty userID", key: "test", value: "test", userID: "", hasError: true},
	}
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()
	s := &PostgresURLStore{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				mock.ExpectQuery(`(?i)insert into urls \(original_url, shorted_url, user_id\) values \(\$1, \$2, \$3\) on conflict \(original_url\) do update set updated_at = now\(\) returning \(created_at = updated_at\) as is_new, shorted_url;`).
					WithArgs(tc.value, tc.key, tc.userID).
					WillReturnRows(sqlmock.NewRows([]string{"is_new", "short_url"}).AddRow(true, "test"))
			}
			_, err := s.SetURL(tc.key, tc.value, tc.userID)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresURLStore_SetBatchURL(t *testing.T) {
	testCases := []struct {
		name      string
		createDTO []models.BatchShortURLCreateDTO
		userID    string
		hasError  bool
	}{
		{name: "Not empty list with correct values", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: "test1"},
			{CorrelationID: "test2", OriginalURL: "https://google.com", ShortURL: "test2"},
		}, userID: "1", hasError: false},
		{name: "Empty list", createDTO: []models.BatchShortURLCreateDTO{}, userID: "1", hasError: false},
		{name: "Not empty list with empty short url", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: ""},
			{CorrelationID: "test2", OriginalURL: "https://google.com", ShortURL: "test2"},
		}, userID: "1", hasError: true},
		{name: "Not empty list with empty original url", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: "test1"},
			{CorrelationID: "test2", OriginalURL: "", ShortURL: "test2"},
		}, userID: "1", hasError: true},
		{name: "Nil list", createDTO: nil, userID: "1", hasError: true},
		{name: "Empty user id", createDTO: nil, userID: "", hasError: true},
	}
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error creating db mock: %v", err)
	}
	defer db.Close()
	s := &PostgresURLStore{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()
			if !tc.hasError {
				for _, dto := range tc.createDTO {
					mock.ExpectQuery(`(?i)insert into urls \(original_url, shorted_url, user_id\) values \(\$1, \$2, \$3\) on conflict \(original_url\) do update set updated_at = now\(\) returning \(created_at = updated_at\) as is_new, shorted_url;`).
						WithArgs(dto.OriginalURL, dto.ShortURL, tc.userID).
						WillReturnRows(sqlmock.NewRows([]string{"is_new", "short_url"}).AddRow(true, "test"))
				}
				mock.ExpectCommit()
			}
			err := s.SetBatchURL(tc.createDTO, tc.userID)
			if tc.hasError {
				assert.NotNil(t, err, "Error is nil")
			} else {
				assert.Nil(t, err, "Error is not nil")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresURLStore_GetURL(t *testing.T) {
	testCases := []struct {
		key    string
		getKey string
		value  string
		name   string
		userID string
	}{
		{name: "Get existing value", key: "test", getKey: "test", value: "test", userID: "1"},
		{name: "Get not existing value", key: "test", getKey: "not exists", value: "test", userID: "1"},
	}
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error create db mock: %v", err)
	}
	defer db.Close()
	s := &PostgresURLStore{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.key == tc.getKey {
				mock.ExpectQuery(`select original_url, deleted_at is not null as is_deleted from urls where shorted_url = \$1`).
					WithArgs(tc.getKey).
					WillReturnRows(sqlmock.NewRows([]string{"original_url", "is_deleted"}).AddRow(tc.value, false))
			} else {
				mock.ExpectQuery(`select original_url, deleted_at is not null as is_deleted from urls where shorted_url = \$1`).
					WithArgs(tc.getKey).
					WillReturnError(sql.ErrNoRows)
			}
			res, err := s.GetURL(tc.getKey)
			if tc.key == tc.getKey {
				assert.Equal(t, res, tc.value, "Get result is not equal to test value")
				assert.Nil(t, err, "Get error is not nil")
			} else {
				assert.Len(t, res, 0, "Get result is not nil")
				assert.NotNil(t, err, "Get error is nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresURLStore_GetUserURLs(t *testing.T) {
	userID := "1"
	testCases := []struct {
		name        string
		getUserID   string
		hasError    bool
		originalURL string
		shortURL    string
	}{
		{name: "Get with correct userID", getUserID: userID, hasError: false, originalURL: "https://ya.ru", shortURL: "test1"},
		{name: "Get with wrong userID", getUserID: "2", hasError: true, originalURL: "https://ya.ru", shortURL: "test1"},
	}
	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error create db mock: %v", err)
	}
	defer db.Close()
	s := &PostgresURLStore{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.hasError {
				mock.ExpectQuery(`select original_url, short_url from urls where user_id = \$1 and deleted_at is null`).
					WithArgs(tc.getUserID).
					WillReturnRows(sqlmock.NewRows([]string{"original_url", "short_url"}).AddRow(tc.originalURL, tc.shortURL))
			} else {
				mock.ExpectQuery(`select original_url, short_url from urls where user_id = \$1 and deleted_at is null`).
					WithArgs(tc.getUserID).
					WillReturnError(sql.ErrNoRows)
			}
			res, err := s.GetUserURLs(tc.getUserID)
			if !tc.hasError {
				assert.Len(t, res, 1, "Get result length is not equal to test value")
				assert.Nil(t, err, "Get error is not nil")
			} else {
				assert.Len(t, res, 0, "Get result is not nil")
				assert.NotNil(t, err, "Get error is nil")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

func TestPostgresURLStore_DeleteURLs(t *testing.T) {
	testCases := []struct {
		name         string
		userID       string
		urls         []string
		hasError     bool
		expectedRows int64
	}{
		{
			name:         "Delete existing URLs",
			userID:       "1",
			urls:         []string{"short1", "short2"},
			hasError:     false,
			expectedRows: 2,
		},
		{
			name:         "Delete non-existent URLs",
			userID:       "1",
			urls:         []string{"short3"},
			hasError:     false,
			expectedRows: 0,
		},
		{
			name:         "Empty URL list",
			userID:       "1",
			urls:         []string{},
			hasError:     true,
			expectedRows: 0,
		},
		{
			name:         "Invalid userID",
			userID:       "",
			urls:         []string{"short1"},
			hasError:     true,
			expectedRows: 0,
		},
	}

	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Errorf("Error creating db mock: %v", err)
	}
	defer db.Close()

	for _, tc := range testCases {
		s := &PostgresURLStore{cfg: &cfg, db: db}
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.urls) > 0 && tc.userID != "" {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)update urls set deleted_at = now\(\) where shorted_url = any\(\$1\) and user_id = \$2 and deleted_at is null`).
					WithArgs(pq.Array(tc.urls), tc.userID).
					WillReturnResult(sqlmock.NewResult(0, tc.expectedRows))
				mock.ExpectCommit()
			}

			err := s.DeleteURLs(tc.userID, tc.urls)
			if tc.hasError {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error during deletion")
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("Not all expectations were met: %v", err)
				}
			}

		})
	}
}

func TestPostgresURLStore_CheckDBConnection(t *testing.T) {
	testCases := []struct {
		hasError bool
		name     string
		error    error
	}{
		{name: "Success", hasError: false, error: nil},
		{name: "Error", hasError: true, error: sql.ErrConnDone},
	}

	cfg := config.GetConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Error create db mock: %v", err)
	}
	defer db.Close()
	s := &PostgresURLStore{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectPing().WillReturnError(tc.error)
			err = s.CheckDBConnection()
			assert.Equal(t, tc.hasError, err != nil, "CheckDBConnection failed")
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}

}
