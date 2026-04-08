package state

import (
	"mime/multipart"

	"github.com/strjkc/noteapi/internal/queries"
	"github.com/strjkc/noteapi/spellcheck"
)

type Storage interface {
	StoreFile(data *multipart.Reader) (string, error)
	GetFile(filePath string) ([]byte, error)
	FileExists(filePath string) bool
	StorageDir() string
	FileURL(fileName string) string
}

type State struct {
	Storage        Storage
	WordMapFactory *spellcheck.WordMapFactory
	DbQueries      *queries.Queries
}

func NewState(storage Storage, wmf *spellcheck.WordMapFactory, dbQueries *queries.Queries) *State {
	s := State{Storage: storage, WordMapFactory: wmf, DbQueries: dbQueries}
	return &s
}
