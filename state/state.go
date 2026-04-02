package state

type Storage interface {
	StoreFile(data []byte, filePath string) (bool, error)
	GetFile(filePath string) ([]byte, error)
	FileExists(filePath string) bool
}

type State struct {
	Storage Storage
}

func NewState(storage Storage) *State {
	s := State{Storage: storage}
	return &s
}
