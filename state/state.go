package state

import (
	"io"
)

type Storage interface {
	StoreFile(data io.ReadCloser, filePath string) (bool, error)
	GetFile(filePath string) ([]byte, error)
	FileExists(filePath string) bool
	GetFilePath(fileName string) string
}

type State struct {
	Storage Storage
}

func NewState(storage Storage) *State {
	s := State{Storage: storage}
	return &s
}
