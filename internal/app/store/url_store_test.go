package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func TestURLStore_SetURL(t *testing.T) {
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
	s := &URLStore{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.key) > 0 && len(tc.value) > 0 {
				mock.ExpectExec(`(?i)insert into urls \(original_url, shorted_url\) values \(\$1, \$2\) on conflict \(original_url\) do update set shorted_url = excluded.shorted_url, updated_at = now\(\);`).
					WithArgs(tc.value, tc.key).
					WillReturnResult(sqlmock.NewResult(1, 1))
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

func TestURLStore_GetURL(t *testing.T) {
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
	s := &URLStore{cfg: &cfg, db: db}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectExec(`(?i)insert into urls \(original_url, shorted_url\) values \(\$1, \$2\) on conflict \(original_url\) do update set shorted_url = excluded.shorted_url, updated_at = now\(\);`).
				WithArgs(tc.value, tc.key).
				WillReturnResult(sqlmock.NewResult(1, 1))
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

func TestURLStore_CheckDBConnection(t *testing.T) {
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
	s := &URLStore{cfg: &cfg, db: db}
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
