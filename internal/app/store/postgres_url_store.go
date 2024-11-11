package store

import (
	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/logger"
)

type PostgresURLStore struct {
	cfg *config.Config
	db  *sql.DB
}

func NewPostgresURLStore(cfg *config.Config) *PostgresURLStore {
	log := logger.NewLogger()
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Log.Error("Error connecting to database", zap.Error(err))
		return nil
	}
	query := `
		select exists (
			select from information_schema.tables 
			where table_schema = 'public' and table_name = 'urls'
		);
	`
	var exists bool
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Log.Error("Error checking if table exists", zap.Error(err))
	}
	if !exists {
		query = `
            create table urls (
                id serial,
				original_url text not null,
                shorted_url text not null,
				created_at timestamp not null default now(),
				updated_at timestamp not null default now(),
				constraint urls_id_pk primary key(id),
				constraint ulrs_original_url_uk unique (original_url)
            );
        `
		_, err = db.Exec(query)
		if err != nil {
			log.Log.Error("Error creating table", zap.Error(err))
		}
	}
	if err != nil {
		log.Log.Error("Error connecting to database", zap.Error(err))
	}
	store := &PostgresURLStore{cfg: cfg, db: db}
	return store
}

func (s *PostgresURLStore) SetURL(key, value string) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if len(value) == 0 {
		return ErrEmptyValue
	}
	log := logger.NewLogger()
	query := `
		insert into urls (original_url, shorted_url) values ($1, $2)
		on conflict (original_url) do update set shorted_url = excluded.shorted_url, updated_at = now();
	`
	_, err := s.db.Exec(query, value, key)
	if err != nil {
		log.Log.Error("Error upserting record", zap.Error(err))
	}
	return nil
}

func (s *PostgresURLStore) GetURL(key string) (string, error) {
	query := `
		select original_url from urls where shorted_url = $1
	`
	var value string
	err := s.db.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return value, nil
}

func (s *PostgresURLStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *PostgresURLStore) CheckDBConnection() error {
	return s.db.Ping()
}
