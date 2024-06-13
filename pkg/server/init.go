package server

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/logging"
	"os"
	"path/filepath"
)

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
