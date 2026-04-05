package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	storageDir string
}

func NewLocalStorage(storagePath string) *LocalStorage {
	l := LocalStorage{storageDir: storagePath}
	return &l
}

func (l *LocalStorage) StoreFile(data io.ReadCloser, fileName string) error {
	if _, err := os.Stat(l.storageDir); os.IsNotExist(err) {
		os.Mkdir(l.storageDir, 0o755)
	}
	filePath := filepath.Join(l.storageDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("error occured")
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, data)
	if err != nil {
		file.Close()
		os.Remove(filePath)
		return err
	}
	return nil
}

func (l *LocalStorage) GetFile(fileName string) ([]byte, error) {
	data, err := os.ReadFile(l.storageDir + fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *LocalStorage) GetStoragePath(fileName string) string {
	if !l.FileExists(fileName) {
		return ""
	}
	return l.storageDir
}

func (l *LocalStorage) FileExists(fileName string) bool {
	if _, err := os.Stat(l.storageDir + fileName); err != nil {
		return false
	}
	return true
}
