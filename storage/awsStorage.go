package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	tmpDir string
	bucket string
	region string
}

func NewS3Storage() *S3Storage {
	tmpDir := os.Getenv("TMPDIR")
	bucket := os.Getenv("BUCKTNAME")
	region := os.Getenv("REGION")
	if tmpDir == "" || bucket == "" || region == "" {
		panic("Unable to initialize storage, and env variable is missing")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Storage{client: client, tmpDir: tmpDir, bucket: bucket, region: region}
}

func (s *S3Storage) StoreFile(data *multipart.Reader) (string, string, error) {
	tmpFile, err := os.CreateTemp(s.tmpDir, "tempFile*")
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
	_, err = tmpFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", "", nil
	}
	contType := "text/markdown"
	obj := s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &fileName,
		Body:        tmpFile,
		ContentType: &contType,
	}
	_, err = s.client.PutObject(context.Background(), &obj)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	return fileName, filepath.Base(tmpFile.Name()), nil
}

func (s *S3Storage) CommitFile(tmpFileName, fileName string) (string, error) {
	tmpPath := path.Join(s.tmpDir, tmpFileName)
	tmpFile, err := os.OpenFile(tmpPath, os.O_RDWR, 777)
	if err != nil {
		return "", err
	}
	contType := "text/markdown"
	obj := s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &fileName,
		Body:        tmpFile,
		ContentType: &contType,
	}
	_, err = s.client.PutObject(context.Background(), &obj)
	if err != nil {
		fmt.Println(err)
	}
	fileURL := s.bucket + fileName
	return fileURL, err
}

func (s *S3Storage) CommitFileWithBackup(tmpFileName, fileName string) (string, string, error) {
	backupName := "backup_" + fileName
	err := s.RenameFile(fileName, "backup_"+fileName)
	if err != nil {
		fmt.Println("An error occured backing up the file")
		return "", "", err
	}
	fileURL, err := s.CommitFile(tmpFileName, fileName)
	if err != nil {
		return "", "", err
	}
	return fileURL, backupName, nil
}

func (s *S3Storage) GetFile(filePath string) ([]byte, error) {
	return nil, nil
}

func (s *S3Storage) FileExists(filePath string) bool {
	params := s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &filePath,
	}
	obj, err := s.client.GetObject(context.Background(), &params)
	if err != nil || obj == nil {
		return false
	}
	return true
}

func (s *S3Storage) StorageDir() string {
	return "https://" + s.bucket + ".amazonaws.com/"
}

func (s *S3Storage) FileURL(fileName string) string {
	if !s.FileExists(fileName) {
		return ""
	}
	return s.bucket + "/" + fileName
}

func (s *S3Storage) DeleteFile(fileName string) error {
	params := s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &fileName,
	}
	_, err := s.client.DeleteObject(context.Background(), &params)
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Storage) RenameFile(fileName, newFileName string) error {
	copySource := s.StorageDir() + fileName
	copyParams := s3.CopyObjectInput{
		Bucket:     &s.bucket,
		CopySource: &copySource,
		Key:        &newFileName,
	}
	_, err := s.client.CopyObject(context.Background(), &copyParams)
	if err != nil {
		return err
	}
	err = s.DeleteFile("mdFile.md")
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Storage) RollbackFileUpdate(backupFileName, fileName string) error {
	err := s.DeleteFile(fileName)
	if err != nil {
		fmt.Println("file exists on disk but we could not delete it")
		return err
	}
	err = s.RenameFile(backupFileName, fileName)
	if err != nil {
		fmt.Println("file exists on disk but we could not delete it")
		return err
	}
	return nil
}
