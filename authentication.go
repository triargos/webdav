package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type UserCredentials struct {
	UserName string
	Realm    string
	Password string
}

func readUserPasswdFile(filePath string) ([]UserCredentials, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %s \n Cause: %s \n", filePath, err.Error())
		return nil, err
	}
	defer CloseFile(file)
	scanner := bufio.NewScanner(file)
	var users []UserCredentials
	for scanner.Scan() {
		line := scanner.Text()
		user, err := readLineToUserCredentials(line)
		if err == nil {
			users = append(users, *user)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func readLineToUserCredentials(line string) (*UserCredentials, error) {
	parts := strings.Split(line, ":")
	if len(parts) == 3 {
		return &UserCredentials{
			UserName: parts[0],
			Realm:    parts[1],
			Password: parts[2],
		}, nil
	}
	return nil, fmt.Errorf("invalid line: %s", line)
}

func verifyCredentials(username string, password string, realm string) bool {
	config, err := readConfiguration()
	if err != nil {
		log.Println("Error reading configuration file")
		return false

	}
	hashToCompare := generateMD5Hash(username, realm, password)
	users, err := readUserPasswdFile(config.CredentialsFilePath)
	if err != nil {
		return false
	}
	for _, user := range users {
		if user.UserName == username && user.Realm == realm {
			return strings.EqualFold(user.Password, hashToCompare)
		}
	}
	return false
}

func generateMD5Hash(userName, realm, password string) string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s:%s:%s", userName, realm, password))
	return fmt.Sprintf("%x", h.Sum(nil))
}
