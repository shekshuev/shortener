package store

import (
	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/shekshuev/shortener/internal/app/logger"
	"github.com/shekshuev/shortener/internal/app/models"
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
				user_id text not null,
				original_url text not null,
                shorted_url text not null,
				created_at timestamp not null default now(),
				updated_at timestamp not null default now(),
				deleted_at timestamp,
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

func (s *PostgresURLStore) SetURL(key, value, userID string) (string, error) {
	if len(key) == 0 {
		return "", ErrEmptyKey
	}
	if len(value) == 0 {
		return "", ErrEmptyValue
	}
	if len(userID) == 0 {
		return "", ErrEmptyUserID
	}
	log := logger.NewLogger()
	query := `
		insert into urls (original_url, shorted_url, user_id) values ($1, $2, $3)
		on conflict (original_url) do update set updated_at = now() 
		returning (created_at = updated_at) as is_new, shorted_url;
	`
	var (
		isNew      bool
		shorterURL string
	)
	err := s.db.QueryRow(query, value, key, userID).Scan(&isNew, &shorterURL)
	if err != nil {
		log.Log.Error("Error upserting record", zap.Error(err))
	}
	if !isNew {
		return shorterURL, ErrAlreadyExists
	}
	return shorterURL, nil
}

func (s *PostgresURLStore) SetBatchURL(createDTO []models.BatchShortURLCreateDTO, userID string) error {
	log := logger.NewLogger()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `
		insert into urls (original_url, shorted_url, user_id) values ($1, $2, $3)
		on conflict (original_url) do update set updated_at = now()
		returning (created_at = updated_at) as is_new, shorted_url;
	`
	hasSameURL := false
	for i := 0; i < len(createDTO); i++ {
		if len(createDTO[i].ShortURL) == 0 {
			return ErrEmptyKey
		}
		if len(createDTO[i].OriginalURL) == 0 {
			return ErrEmptyValue
		}
		if len(userID) == 0 {
			return ErrEmptyUserID
		}
		var (
			isNew    bool
			shortURL string
		)
		err := s.db.QueryRow(query, createDTO[i].OriginalURL, createDTO[i].ShortURL, userID).Scan(&isNew, &shortURL)
		if err != nil {
			log.Log.Error("Error upserting record", zap.Error(err))
			tx.Rollback()
		}
		if !isNew {
			createDTO[i].ShortURL = shortURL
			hasSameURL = true
		}
	}
	if hasSameURL {
		return ErrAlreadyExists
	}
	return tx.Commit()
}

func (s *PostgresURLStore) GetURL(key, userID string) (string, error) {
	query := `
		select original_url, deleted_at is not null as is_deleted from urls where shorted_url = $1 and user_id = $2;
	`
	var value string
	var isDeleted bool
	err := s.db.QueryRow(query, key, userID).Scan(&value, &isDeleted)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	if isDeleted {
		return "", ErrAlreadyDeleted
	}
	return value, nil
}

func (s *PostgresURLStore) GetUserURLs(userID string) ([]models.UserShortURLReadDTO, error) {
	query := `
		select original_url, short_url from urls where user_id = $1 and deleted_at is null;
	`
	var readDTO []models.UserShortURLReadDTO
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, ErrNotFound
	}
	for rows.Next() {
		var (
			originalURL string
			shortURL    string
		)
		err := rows.Scan(&originalURL, &shortURL)
		if err != nil {
			return nil, err
		}
		readDTO = append(readDTO, models.UserShortURLReadDTO{ShortURL: shortURL, OriginalURL: originalURL})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return readDTO, nil
}

func (s *PostgresURLStore) DeleteURLs(userID string, urls []string) error {
	if len(userID) == 0 {
		return ErrEmptyUserID
	}
	if len(urls) == 0 {
		return ErrEmptyURLs
	}

	ch := make(chan string, len(urls))
	for _, url := range urls {
		ch <- url
	}
	close(ch)

	const workers = 4
	results := make(chan []string, workers)
	for i := 0; i < workers; i++ {
		go func() {
			var batch []string
			for url := range ch {
				batch = append(batch, url)
				if len(batch) >= 100 {
					results <- batch
					batch = nil
				}
			}
			if len(batch) > 0 {
				results <- batch
			}
		}()
	}

	go func() {
		for i := 0; i < workers; i++ {
			results <- nil
		}
		close(results)
	}()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
        update urls set deleted_at = now() where shorted_url = any($1) and user_id = $2 and deleted_at is null;
    `
	for batch := range results {
		if batch == nil {
			continue
		}
		_, err = tx.Exec(query, pq.Array(batch), userID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
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
