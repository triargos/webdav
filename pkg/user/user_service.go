package user

import (
	"errors"
	"fmt"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/fs"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Service interface {
	AddUser(username string, user config.User) error
	GetUser(username string) config.User
	GetUsers() map[string]config.User
	HasUser(username string) bool
	RemoveUser(username string) error
	InitializeDirectories() error
	HashPasswords() error
}

type OsUserService struct {
	configService config.Service
	fsService     fs.Service
}

func NewOsUserService(configService config.Service, fsService fs.Service) Service {
	return &OsUserService{configService: configService, fsService: fsService}
}

func (s *OsUserService) createUserDirectories(user config.User, contentRoot string) error {
	userRoot := filepath.Join(contentRoot, user.Root)
	createDirectoryErr := s.fsService.CreateDirectories(userRoot, os.ModePerm)
	if createDirectoryErr != nil {
		return fmt.Errorf("failed to create users root directory: %v", createDirectoryErr)
	}
	for _, subdirectory := range user.SubDirectories {
		subDirectoryPath := filepath.Join(userRoot, subdirectory)
		createSubDirectoryErr := s.fsService.CreateDirectories(subDirectoryPath, os.ModePerm)
		if createSubDirectoryErr != nil {
			slog.Error("failed to create subdirectory", "subdirectory", subdirectory, "error", createSubDirectoryErr)
		}
	}
	return nil
}

func (s *OsUserService) AddUser(username string, user config.User) error {
	hash := GenHash([]byte(user.Password))
	user.Password = hash
	s.configService.AddUser(username, user)
	writeConfigErr := s.configService.Write()
	if writeConfigErr != nil {
		return fmt.Errorf("failed to write config file: %s", writeConfigErr)
	}
	contentRoot := s.configService.Get().Content.Dir
	createUserDirectoriesErr := s.createUserDirectories(user, contentRoot)
	if createUserDirectoriesErr != nil {
		return createUserDirectoriesErr
	}
	return nil
}

func (s *OsUserService) GetUser(username string) config.User {
	users := s.configService.Get().Users
	return users[username]
}

func (s *OsUserService) HasUser(username string) bool {
	users := s.configService.Get().Users
	_, ok := users[username]
	return ok
}

func (s *OsUserService) RemoveUser(username string) error {
	if !s.HasUser(username) {
		return errors.New("user does not exist")
	}
	user := s.GetUser(username)
	dirPath := filepath.Join(s.configService.Get().Content.Dir, user.Root)
	removeUserDirectoryErr := s.fsService.RemoveDirectories(dirPath)
	if removeUserDirectoryErr != nil {
		return fmt.Errorf("error removing user directory: %s", removeUserDirectoryErr)
	}
	s.configService.RemoveUser(username)
	writeConfigErr := s.configService.Write()
	if writeConfigErr != nil {
		return fmt.Errorf("failed to write config file: %s", writeConfigErr)
	}
	return nil
}

func (s *OsUserService) InitializeDirectories() error {
	users := s.configService.Get().Users
	contentRoot := s.configService.Get().Content.Dir
	for _, user := range users {
		createDirectoriesErr := s.createUserDirectories(user, contentRoot)
		if createDirectoriesErr != nil {
			slog.Error("failed to initialize directories", "error", createDirectoriesErr)
		}
	}
	return nil
}

func (s *OsUserService) HashPasswords() error {
	for username, user := range s.configService.Get().Users {
		if !isHashed(user.Password) {
			slog.Info("password for user is not hashed, hashing now", "username", username)
			user.Password = GenHash([]byte(user.Password))
			s.configService.UpdateUser(username, user)
		}
	}
	return s.configService.Write()
}

func (s *OsUserService) GetUsers() map[string]config.User {
	return s.configService.Get().Users
}

func isHashed(password string) bool {
	// bcrypt hashed passwords have a specific format that starts with "$2a$", "$2b$", or "$2y$"
	return strings.HasPrefix(password, "$2a$") || strings.HasPrefix(password, "$2b$") || strings.HasPrefix(password, "$2y$")
}

func GenHash(password []byte) string {
	pw, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		log.Fatal(err)
	}
	return string(pw)
}
