package main

import (
	"encoding/json"
	"os"
	"strings"
)

func CheckPathPermission(path string, username string) bool {
	config, err := readConfiguration()
	if err != nil {
		return handleConfigurationReadError(err)
	}
	if username == config.AdminUserName {
		return true
	}
	relevantPath := strings.ToLower(strings.Split(path, "/")[1])
	if relevantPath == config.UsersRoot {
		return isUserFolder(username, path)
	}
	restricted, requiredPermission := isRestricted(relevantPath, config)
	if restricted {
		permissions, err := readPermissions(config.PermissionsFilePath)
		if err != nil {
			return handleConfigurationReadError(err)
		}
		return hasPermission(username, requiredPermission, permissions)
	}
	return true
}

func isUserFolder(username string, path string) bool {
	usernameSegment := strings.Split(path, "/")[2]
	return usernameSegment == username
}

type Permissions struct {
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
}

func readPermissions(filePath string) ([]Permissions, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer CloseFile(file)
	decoder := json.NewDecoder(file)
	var permissions []Permissions
	err = decoder.Decode(&permissions)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func isRestricted(path string, config *Configuration) (bool, string) {
	for _, restrictedDir := range config.RestrictedDirectories {
		if path == restrictedDir {
			return true, restrictedDir
		}
	}
	return false, ""
}

func hasPermission(username string, requiredPermission string, permissions []Permissions) bool {
	for _, permissionObject := range permissions {
		if permissionObject.Username == username {
			for _, permission := range permissionObject.Permissions {
				if permission == requiredPermission {
					return true
				}
			}
		}
	}
	return false
}
