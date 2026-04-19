package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
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
	return fileName, filepath.Base(tmpFile.Name()), nil
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

func (l *LocalStorage) CommitFile(tmpFileName, fileName string) (string, error) {
	err := l.RenameFile(tmpFileName, fileName)
	if err != nil {
		fmt.Println("An error occured commiting the file")
		return "", err
	}
	fileURL := path.Join(l.storageDir, fileName)
	return fileURL, nil
}

func (l *LocalStorage) CommitFileWithBackup(tmpFileName, fileName string) (string, string, error) {
	backupName := "backup_" + fileName
	err := l.RenameFile(fileName, backupName)
	if err != nil {
		fmt.Println("An error occured creating a file backup")
		return "", "", err
	}
	fileURL, err := l.CommitFile(tmpFileName, fileName)
	if err != nil {
		return "", "", err
	}
	return backupName, fileURL, nil
}

func (l *LocalStorage) RollbackFileUpdate(backupFileName, fileName string) error {
	err := l.DeleteFile(fileName)
	if err != nil {
		fmt.Println("file exists on disk but we could not delete it")
		return err
	}
	err = l.RenameFile(backupFileName, fileName)
	if err != nil {
		fmt.Println("file exists on disk but we could not delete it")
		return err
	}
	return err
}
