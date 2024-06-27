package server

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/fs"
	"log/slog"
	"os"
	"path/filepath"
)

func CreateUserDirectories() {
	cfg := config.Get()
	for _, user := range *cfg.Users {
		rootPath := filepath.Join(cfg.Content.Dir, user.Root)
		if !fs.PathExists(rootPath) {
			err := os.MkdirAll(rootPath, os.ModePerm)
			if err != nil {
				slog.Error("Failed to create root directory for user", "path", user.Root, "error", err)
			}
		}
		subDirectories := user.SubDirectories
		for _, dir := range subDirectories {
			subDirectoryPath := filepath.Join(rootPath, dir)
			if !fs.PathExists(subDirectoryPath) {
				err := os.MkdirAll(subDirectoryPath, os.ModePerm)
				if err != nil {
					slog.Error("Failed to create subdirectory for user", "path", dir, "error", err)
				}
			}
		}
	}
}
