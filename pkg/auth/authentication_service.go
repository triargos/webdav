package auth

import (
	"github.com/triargos/webdav/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strings"
)

type Service interface {
	Authenticate(username, password string) bool
	HasPermission(path string, username string) bool
}

type ServiceImpl struct {
	userService user.Service
}

func New(userService user.Service) Service {
	return &ServiceImpl{userService: userService}
}

func (s *ServiceImpl) Authenticate(username, password string) bool {
	if !s.userService.HasUser(username) {
		return false
	}
	userObject := s.userService.GetUser(username)
	verifyPasswordErr := bcrypt.CompareHashAndPassword([]byte(userObject.Password), []byte(password))
	return verifyPasswordErr == nil
}

func (s *ServiceImpl) HasPermission(path string, username string) bool {
	slog.Info("User has permission", "path", path, "username", username)
	if !s.userService.HasUser(username) {
		return false
	}
	userObject := s.userService.GetUser(username)
	if userObject.Admin {
		return true
	}
	if userObject.Jail {
		return isSubPath(userObject.Root, path)
	}
	for otherUsername, otherUser := range s.userService.GetUsers() {
		if otherUsername == username {
			continue
		}
		if otherUser.Root != "" && isSubPath(otherUser.Root, path) {
			return false
		}
	}
	return true

}

func isSubPath(parent, child string) bool {
	parent = strings.TrimSuffix(parent, "/") + "/"
	child = strings.TrimSuffix(child, "/") + "/"
	return strings.HasPrefix(child, parent)
}
