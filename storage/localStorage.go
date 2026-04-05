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

func (l *LocalStorage) StoreFile(data io.ReadCloser, fileName string) (bool, error) {
	// TODO: what if file already exits
	fmt.Printf("\n%s\n", l.storageDir)
	if _, err := os.Stat(l.storageDir); os.IsNotExist(err) {
		os.Mkdir(l.storageDir, 0o755)
	}
	path := filepath.Join(l.storageDir, fileName)
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("error occured")
		return false, err
	}
	defer file.Close()
	io.Copy(file, data)
	return true, nil
}

func (l *LocalStorage) GetFile(fileName string) ([]byte, error) {
	data, err := os.ReadFile(l.storageDir + fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *LocalStorage) GetFilePath(fileName string) string {
	if _, err := os.Stat(l.storageDir + fileName); err != nil {
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
