package auth

import (
	"github.com/triargos/webdav/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Service interface {
	Authenticate(username, password string) bool
	HasPermission(path string, username string) bool
}

type BasicAuthenticator struct {
	userService user.Service
}

func New(userService user.Service) Service {
	return &BasicAuthenticator{userService: userService}
}

func (s *BasicAuthenticator) Authenticate(username, password string) bool {
	if !s.userService.HasUser(username) {
		return false
	}
	userObject := s.userService.GetUser(username)
	verifyPasswordErr := bcrypt.CompareHashAndPassword([]byte(userObject.Password), []byte(password))
	return verifyPasswordErr == nil
}

func (s *BasicAuthenticator) HasPermission(path string, username string) bool {
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
	return strings.HasPrefix(strings.ToLower(child), strings.ToLower(parent))
}
