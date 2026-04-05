package storage

/*
import (
	"os"
)

type LocalStorage struct {
	storageDir string
}

func NewLocalStorage(storagePath string) *LocalStorage {
	l := LocalStorage{storageDir: storagePath}
	return &l
}

func (l *LocalStorage) StoreFile(data []byte, fileName string) (bool, error) {
	err := os.WriteFile(l.storageDir+fileName, data, 0o777)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *LocalStorage) GetFile(fileName string) ([]byte, error) {
	data, err := os.ReadFile(l.storageDir + fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *LocalStorage) FileExists(fileName string) bool {
	if _, err := os.Stat(l.storageDir + fileName); err != nil {
		return false
	}
	return true
}
*/
