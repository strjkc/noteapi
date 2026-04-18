package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	tmp    string
}

func NewS3Storage() *S3Storage {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	return &S3Storage{client: client, tmp: "/temp/"}
}

func (s *S3Storage) StoreFile(data *multipart.Reader) (string, string, error) {
	tmpFile, err := os.CreateTemp(s.tmp, "tempFile*")
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
	bucket := "noteapi-files426871656705-us-east-1-an"
	key := "somekey"
	contType := "text/markdown"
	obj := s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        tmpFile,
		ContentType: &contType,
	}
	s.client.PutObject(context.Background(), &obj)
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
