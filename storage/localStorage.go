package storage

func StoreFile(data []byte, filePath string) (bool, err) {}
func GetFile(filePath string) ([]byte, error)            {}
func FileExists(filePath string) bool                    {}
