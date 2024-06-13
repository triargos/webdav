package server

import (
	"github.com/triargos/webdav/pkg/config"
)

func AuthenticateUser(username, password string) bool {
	//Find the user
	user := config.Value.Users[username]
	isValid := user.Password == password
	return isValid
}
