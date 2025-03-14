package store

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/shekshuev/shortener/internal/app/models"
)

// CreateSnapshot создаёт снапшот хранилища в файл.
func (s *MemoryURLStore) CreateSnapshot() error {
	file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range s.urls {
		urlData := models.SerializeData{
			UserID:      value.UserID,
			ShortURL:    key,
			OriginalURL: value.URL,
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

// LoadSnapshot загружает данные из файла снапшота в хранилище.
func (s *MemoryURLStore) LoadSnapshot() error {
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
		s.urls[urlData.ShortURL] = UserURL{UserID: urlData.UserID, URL: urlData.OriginalURL}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
