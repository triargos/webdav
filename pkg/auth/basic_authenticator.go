package auth

import (
	"fmt"
	"github.com/triargos/webdav/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type BasicAuthenticator struct {
	userService user.Service
}

func (authenticator BasicAuthenticator) PerformAuthentication(writer http.ResponseWriter, request *http.Request) (string, http.ResponseWriter) {
	username, password, ok := request.BasicAuth()
	if !ok {
		writer.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
		return "", writer
	}
	validateCredentialsErr := authenticator.validateCredentials(username, password)
	if validateCredentialsErr != nil {
		return "", writer
	}
	return username, writer
}

func (authenticator BasicAuthenticator) validateCredentials(username, password string) error {
	if !authenticator.userService.HasUser(username) {
		return fmt.Errorf("user not found")
	}
	userObject := authenticator.userService.GetUser(username)
	verifyPasswordErr := bcrypt.CompareHashAndPassword([]byte(userObject.Password), []byte(password))
	return verifyPasswordErr
}

func NewBasicAuthenticator(userService user.Service) Authenticator {
	return &BasicAuthenticator{userService: userService}
}
