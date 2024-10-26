package store

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/shekshuev/shortener/internal/app/models"
)

func (s *URLStore) CreateSnapshot() error {
	file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range s.urls {
		urlData := models.SerializeData{
			UUID:        uuid.New().String(),
			ShortURL:    key,
			OriginalURL: value,
		}

		data, err := json.Marshal(urlData)
		if err != nil {
			return err
		}

		_, err = file.Write(append(data, '\n'))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *URLStore) LoadSnapshot() error {
	if _, err := os.Stat(s.cfg.FileStoragePath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(s.cfg.FileStoragePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var urlData models.SerializeData
		err := json.Unmarshal(scanner.Bytes(), &urlData)
		if err != nil {
			continue
		}
		s.urls[urlData.ShortURL] = urlData.OriginalURL
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
