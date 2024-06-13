package server

import (
	"github.com/triargos/webdav/pkg/config"
	"log"
	"strings"
)

// CheckPermission verifies if the user has access to the specified path.
func CheckPermission(currentPath, username string) bool {
	user, ok := config.Value.Users[username]
	if !ok {
		// User not found in configuration
		return false
	}

	// Admin users have access to all paths
	if user.Admin {
		return true
	}

	// Jail users to their root directory
	if user.Jail {
		// Check if the currentPath is a subpath of the user's root directory
		return isSubPath(user.Root, currentPath)
	}
	log.Println("Checking permission for user", username, "on path", currentPath)

	// Check if the currentPath is within another user's directory
	for otherUserName, otherUser := range config.Value.Users {
		if otherUserName == username {
			continue
		}
		// If the path is within another user's root, deny access
		if otherUser.Root != "" && isSubPath(otherUser.Root, currentPath) {
			return false
		}
	}
	log.Println("User", username, "has permission to access", currentPath)

	// Default to allowing access if no other conditions block it
	return true
}

// isSubPath checks if the child path is a subpath of the parent path.
func isSubPath(parent, child string) bool {
	// Ensure both paths are normalized
	parent = strings.TrimSuffix(parent, "/") + "/"
	child = strings.TrimSuffix(child, "/") + "/"
	return strings.HasPrefix(child, parent)
}
