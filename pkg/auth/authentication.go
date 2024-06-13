package auth

import (
	"github.com/triargos/webdav/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(username, password string) bool {
	cfg := config.Get()
	user := (*cfg.Users)[username]
	isValid := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return isValid == nil
}
