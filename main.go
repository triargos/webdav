package main

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/logging"
	"github.com/triargos/webdav/pkg/server"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if !fs.PathExists("/etc/webdav") {
		err := os.Mkdir("/etc/webdav", 0755)
		if err != nil {
			log.Fatalf("Error creating /etc/webdav directory: %s\n", err)
		}
	}
	logging.InitLoggers()
	logging.Log.Info.Println("Logging initialized")
	logging.Log.Info.Println("Reading config...")
	err := config.Init()
	if err != nil {
		logging.Log.Error.Fatalf("Error reading config: %s\n", err)
	}
	logging.Log.Info.Println("Config read successfully")
	logging.Log.Info.Println("Creating user directories...")
	CreateUserDirectories()
	logging.Log.Info.Println("User directories created successfully")
	logging.Log.Info.Println("Hashing non-hashed passwords...")
	server.HashPasswords()
	logging.Log.Info.Println("Passwords hashed successfully")
	logging.Log.Info.Println("Starting server...")
	err = server.StartWebdavServer()
	if err != nil {
		logging.Log.Error.Fatalf("Error starting server: %s\n", err)
	}
}

func CreateUserDirectories() {
	for _, user := range *config.Value.Users {
		rootPath := filepath.Join(config.Value.Content.Dir, user.Root)
		if !fs.PathExists(rootPath) {
			err := os.MkdirAll(rootPath, os.ModePerm)
			if err != nil {
				logging.Log.Error.Printf("Error creating user root directory: %s\n", err)
			}
		}
		subDirectories := user.SubDirectories
		for _, dir := range subDirectories {
			subDirectoryPath := filepath.Join(rootPath, dir)
			if !fs.PathExists(subDirectoryPath) {
				err := os.MkdirAll(subDirectoryPath, os.ModePerm)
				if err != nil {
					logging.Log.Error.Printf("Error creating user subdirectory: %s\n", err)
				}
			}
		}
	}
}
