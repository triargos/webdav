package fs

import "os"

type Service interface {
	CreateDirectories(path string, mode os.FileMode) error
	RemoveDirectories(path string) error
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
