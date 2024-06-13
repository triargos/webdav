package server

import (
	"github.com/triargos/webdav/pkg/config"
)

func AuthenticateUser(username, password string) bool {
	user := config.Value.Users[username]
	isValid := user.Password != "" && user.Password == password
	return isValid
}
