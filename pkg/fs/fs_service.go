package fs

import "os"

type Service interface {
	CreateDirectories(path string, mode os.FileMode) error
	RemoveDirectories(path string) error

	ReadFileContent(path string) ([]byte, error)
	ReadFile(path string) (*os.File, error)

	WriteFileContent(path string, content []byte, mode os.FileMode) error
	CreateFile(path string) (*os.File, error)
}

type OsFileSystemService struct{}

func NewOsFileSystemService() Service {
	return OsFileSystemService{}
}

func (s OsFileSystemService) CreateDirectories(path string, mode os.FileMode) error {
	return os.MkdirAll(path, mode)
}

func (s OsFileSystemService) RemoveDirectories(path string) error {
	return os.RemoveAll(path)
}

func (s OsFileSystemService) ReadFileContent(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (s OsFileSystemService) ReadFile(path string) (*os.File, error) {
	return os.Open(path)
}

func (s OsFileSystemService) WriteFileContent(path string, content []byte, mode os.FileMode) error {
	return os.WriteFile(path, content, mode)
}

func (s OsFileSystemService) CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}
