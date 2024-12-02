package store

import (
	"os"
	"testing"

	"github.com/shekshuev/shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func removeTestFile(filePath string) {
	_ = os.Remove(filePath)
}

func TestCreateSnapshot(t *testing.T) {
	cfg := config.GetConfig()
	removeTestFile(cfg.FileStoragePath)

	store := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}

	store.urls["short1"] = UserURL{UserID: "1", URL: "https://ya.ru"}
	store.urls["short2"] = UserURL{UserID: "1", URL: "https://google.com"}

	err := store.CreateSnapshot()
	assert.Nil(t, err, "Error should be nil when creating snapshot")

	fileInfo, err := os.Stat(cfg.FileStoragePath)
	assert.Nil(t, err, "Error should be nil when checking file stats")
	assert.False(t, fileInfo.IsDir(), "Snapshot should create a valid file")
	removeTestFile(cfg.FileStoragePath)
}

func TestLoadSnapshot(t *testing.T) {
	cfg := config.GetConfig()
	removeTestFile(cfg.FileStoragePath)

	store := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}

	store.urls["short1"] = UserURL{UserID: "1", URL: "https://ya.ru"}
	store.urls["short2"] = UserURL{UserID: "1", URL: "https://google.com"}
	err := store.CreateSnapshot()
	assert.Nil(t, err, "Error should be nil when creating snapshot")

	store2 := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}
	err = store2.LoadSnapshot()
	assert.Nil(t, err, "Error should be nil when loading snapshot")

	assert.Equal(t, store.urls["short1"], store2.urls["short1"], "Loaded value for short1 does not match")
	assert.Equal(t, store.urls["short2"], store2.urls["short2"], "Loaded value for short2 does not match")
	removeTestFile(cfg.FileStoragePath)
}

func TestLoadSnapshot_FileDoesNotExist(t *testing.T) {
	cfg := config.GetConfig()

	store := &MemoryURLStore{urls: make(map[string]UserURL), cfg: &cfg}

	err := store.LoadSnapshot()
	assert.Nil(t, err, "Error should be nil when loading snapshot from non-existent file")
	assert.Empty(t, store.urls, "URL map should be empty when loading from a non-existent file")
}
