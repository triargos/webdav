package main

import (
	"encoding/json"
	"errors"
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
	accessControl, err := readAccessControl(config.PermissionsFilePath)
	if err != nil {
		return handleConfigurationReadError(err)

	}
	permissions, err := findUserPermissions(username, *accessControl)
	if err != nil {
		return false
	}
	relevantPath := strings.ToLower(strings.Split(path, "/")[1])
	if relevantPath == config.UsersRoot {
		return isUserFolder(username, path)
	}
	//Check if the path is restricted for the user specifically
	restrictedForUser := findRestriction(relevantPath, permissions.Restricted)
	if restrictedForUser {
		return false
	}
	//Check if the path is restricted for all
	restrictedForAll := findRestriction(relevantPath, accessControl.RestrictedDirectories)
	if !restrictedForAll {
		return true
	}
	//Otherwise, check if the user has the required permission
	return hasPermission(relevantPath, permissions.Allowed)
}

func isUserFolder(username string, path string) bool {
	usernameSegment := strings.Split(path, "/")[2]
	return usernameSegment == username
}

type AccessControl struct {
	RestrictedDirectories []string          `json:"restrictedDirectories"`
	UserPermissions       []UserPermissions `json:"userPermissions"`
}

type UserPermissions struct {
	Username   string   `json:"username"`
	Allowed    []string `json:"allowed"`
	Restricted []string `json:"restricted"`
}

func readAccessControl(filePath string) (*AccessControl, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer CloseFile(file)
	decoder := json.NewDecoder(file)
	var permissions AccessControl
	err = decoder.Decode(&permissions)
	if err != nil {
		return nil, err
	}
	return &permissions, nil
}

func findUserPermissions(username string, accessControl AccessControl) (*UserPermissions, error) {
	for _, userPermission := range accessControl.UserPermissions {
		if userPermission.Username == username {
			return &userPermission, nil
		}
	}
	return nil, errors.New("user not found")
}

func findRestriction(path string, paths []string) bool {
	for _, pathElement := range paths {
		if pathElement == path {
			return true
		}
	}
	return false
}

func hasPermission(requiredPermission string, allowed []string) bool {
	for _, permission := range allowed {
		if permission == requiredPermission {
			return true
		}
	}
	return false
}
