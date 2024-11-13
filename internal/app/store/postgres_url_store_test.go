package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestPostgresURLStore_SetURL(t *testing.T) {
	testCases := []struct {
		key   string
		value string
		name  string
	}{
		{name: "Normal key and value", key: "test", value: "test"},
		{name: "Empty key", key: "", value: "test"},
		{name: "Empty value", key: "test", value: ""},
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
			if len(tc.key) > 0 && len(tc.value) > 0 {
				mock.ExpectQuery(`(?i)insert into urls \(original_url, shorted_url\) values \(\$1, \$2\) on conflict \(original_url\) do update set shorted_url = excluded.shorted_url, updated_at = now\(\) returning \(created_at = updated_at\) as is_new;`).
					WithArgs(tc.value, tc.key).
					WillReturnRows(sqlmock.NewRows([]string{"is_new"}).AddRow(true))
			}
			err := s.SetURL(tc.key, tc.value)
			if len(tc.key) == 0 || len(tc.value) == 0 {
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
		hasError  bool
	}{
		{name: "Not empty list with correct values", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: "test1"},
			{CorrelationID: "test2", OriginalURL: "https://google.com", ShortURL: "test2"},
		}, hasError: false},
		{name: "Empty list", createDTO: []models.BatchShortURLCreateDTO{}, hasError: false},
		{name: "Not empty list with empty short url", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: ""},
			{CorrelationID: "test2", OriginalURL: "https://google.com", ShortURL: "test2"},
		}, hasError: true},
		{name: "Not empty list with empty original url", createDTO: []models.BatchShortURLCreateDTO{
			{CorrelationID: "test1", OriginalURL: "https://ya.ru", ShortURL: "test1"},
			{CorrelationID: "test2", OriginalURL: "", ShortURL: "test2"},
		}, hasError: true},
		{name: "Nil list", createDTO: nil, hasError: true},
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
					mock.ExpectQuery(`(?i)insert into urls \(original_url, shorted_url\) values \(\$1, \$2\) on conflict \(original_url\) do update set shorted_url = excluded.shorted_url, updated_at = now\(\) returning \(created_at = updated_at\) as is_new;`).
						WithArgs(dto.OriginalURL, dto.ShortURL).
						WillReturnRows(sqlmock.NewRows([]string{"is_new"}).AddRow(true))
				}
				mock.ExpectCommit()
			}
			err := s.SetBatchURL(tc.createDTO)
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
	}{
		{name: "Get existing value", key: "test", getKey: "test", value: "test"},
		{name: "Get not existing value", key: "test", getKey: "not exists", value: "test"},
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
			mock.ExpectQuery(`(?i)insert into urls \(original_url, shorted_url\) values \(\$1, \$2\) on conflict \(original_url\) do update set shorted_url = excluded.shorted_url, updated_at = now\(\) returning \(created_at = updated_at\) as is_new;`).
				WithArgs(tc.value, tc.key).
				WillReturnRows(sqlmock.NewRows([]string{"is_new"}).AddRow(true))
			err := s.SetURL(tc.key, tc.value)
			assert.Nil(t, err, "Set error is not nil")
			if tc.key == tc.getKey {
				mock.ExpectQuery(`select original_url from urls where shorted_url = \$1`).
					WithArgs(tc.getKey).
					WillReturnRows(sqlmock.NewRows([]string{"original_url"}).AddRow(tc.value))
			} else {
				mock.ExpectQuery(`select original_url from urls where shorted_url = \$1`).
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
