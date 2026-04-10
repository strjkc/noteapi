package storage

import (
	"fmt"
	"io"
	"mime/multipart"
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

func (l *LocalStorage) StoreFile(data *multipart.Reader) (string, string, error) {
	if _, err := os.Stat(l.storageDir); os.IsNotExist(err) {
		os.Mkdir(l.storageDir, 0o755)
	}
	tmpFile, err := os.CreateTemp(l.storageDir, "tempFile*")
	if err != nil {
		fmt.Println("error occured")
		return "", "", err
	}
	defer tmpFile.Close()
	// defer os.Remove(tmpFile.Name())

	var fileName string

	for {
		part, err := data.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", "", err
		}
		if part.FileName() == "" {
			continue
		}
		if fileName == "" {
			fileName = part.FileName()
		}

		_, err = io.Copy(tmpFile, part)
		if err != nil {
			return "", "", err
		}
	}

	/*
		filePath := filepath.Join(l.storageDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
		tmpFile.Seek(0, io.SeekStart)
		_, err = io.Copy(file, tmpFile)
		if err != nil {
			os.Remove(filePath)
			return "", err
		}
	*/
	return fileName, tmpFile.Name(), nil
}

func (l *LocalStorage) GetFile(fileName string) ([]byte, error) {
	data, err := os.ReadFile(l.storageDir + fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *LocalStorage) StorageDir() string {
	return l.storageDir
}

func (l *LocalStorage) FileURL(fileName string) string {
	if !l.FileExists(fileName) {
		return ""
	}
	return filepath.Join(l.storageDir, fileName)
}

func (l *LocalStorage) FileExists(fileName string) bool {
	if _, err := os.Stat(l.storageDir + fileName); err != nil {
		return false
	}
	return true
}

func (l *LocalStorage) DeleteFile(fileName string) error {
	err := os.Remove(l.storageDir + fileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File: %s not found, in dir: %s. Nothhing was changed\n", fileName, l.storageDir)
		}
		return err
	}
	return nil
}

func (l *LocalStorage) RenameFile(fileName, newFileName string) error {
	oldPath := filepath.Join(l.storageDir, fileName)
	newPath := filepath.Join(l.storageDir, newFileName)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	return nil
}
