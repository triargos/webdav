package auth

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/logging"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func HashPasswords() {
	users := *config.Value.Users
	for username, user := range users {
		if !isHashed(user.Password) {
			logging.Log.Info.Println("Hashing password for user", username)
			user.Password = GenHash([]byte(user.Password))
			users[username] = user
		}
	}
	config.Value.Users = &users
	err := config.WriteConfig(config.Value)
	if err != nil {
		log.Fatal(err)
	}
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
