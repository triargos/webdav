package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func CheckPathPermission(path string, username string) (bool, error) {
	configuration, err := readConfiguration()
	if err != nil {
		return false, err
	}
	accessControl, err := readAccessControl(configuration.PermissionsFilePath)
	if err != nil {
		return false, err
	}
	isRestricted := isRestrictedDir(path, accessControl)
	if !isRestricted {
		return true, nil
	}
	userPermissions, err := findUserPermissions(username, accessControl)
	if err != nil {
		return false, err
	}
	//Check if the path is explitly forbidden
	if isForbidden(path, userPermissions.Restricted) {
		return true, nil
	}
	return hasPermission(path, userPermissions.Allowed), nil
}

func isRestrictedDir(path string, accessControl *AccessControl) bool {
	for _, dir := range accessControl.RestrictedDirectories {
		if strings.HasPrefix(path, dir) {
			return true
		}
	}
	return false
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

func findUserPermissions(username string, accessControl *AccessControl) (*UserPermissions, error) {
	for _, userPermission := range accessControl.UserPermissions {
		if userPermission.Username == username {
			return &userPermission, nil
		}
	}
	return nil, errors.New("user not found")
}

// TODO: Fix this
func isForbidden(path string, forbidden []string) bool {
	for _, forbiddenPath := range forbidden {
		if strings.HasPrefix(path, forbiddenPath) {
			println("Forbidden:", forbiddenPath)
			return true
		}
	}
	return false
}

func hasPermission(path string, allowed []string) bool {
	for _, allowedPath := range allowed {
		if strings.HasPrefix(path, allowedPath) {
			return true
		}
	}
	return false
}
