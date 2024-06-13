package auth

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(username, password string) bool {
	logging.Log.Info.Println("Authenticating user", username)
	users := *config.Value.Users
	user := users[username]
	isValid := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return isValid == nil
}
