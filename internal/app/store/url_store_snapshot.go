package store

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/shekshuev/shortener/internal/app/models"
)

func CreateSnapshot(store *URLStore) error {
	file, err := os.OpenFile(store.cfg.FileStoragePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range store.urls {
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

func LoadSnapshot(store *URLStore) error {
	if _, err := os.Stat(store.cfg.FileStoragePath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(store.cfg.FileStoragePath)
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
		store.urls[urlData.ShortURL] = urlData.OriginalURL
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
