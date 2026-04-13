package storage

import "mime/multipart"

type S3Storage struct {
	storageDir string
}

func NewS3Storage(storagePath string) *S3Storage {
	s3 := S3Storage{storageDir: storagePath}
	return &s3
}

func (s3 *S3Storage) StoreFile(data *multipart.Reader) (string, string, error) {
	return "", "", nil
}

func (s3 *S3Storage) GetFile(filePath string) ([]byte, error) {
	return nil, nil
}

func (s3 *S3Storage) FileExists(filePath string) bool {
	return false
}

func (s3 *S3Storage) StorageDir() string {
	return ""
}

func (s3 *S3Storage) FileURL(fileName string) string {
	return ""
}

func (s3 *S3Storage) DeleteFile(fileName string) error {
	return nil
}

func (s3 *S3Storage) RenameFile(fileName, newFileName string) error {
	return nil
}
