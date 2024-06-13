package server

import (
	"github.com/triargos/webdav/pkg/config"
	"strings"
)

func CheckPermission(currentPath, username string) bool {
	users := *config.Value.Users
	user, ok := users[username]
	if !ok {
		return false
	}
	if user.Admin {
		return true
	}
	if user.Jail {
		return isSubPath(user.Root, currentPath)
	}
	for otherUserName, otherUser := range *config.Value.Users {
		if otherUserName == username {
			continue
		}
		if otherUser.Root != "" && isSubPath(otherUser.Root, currentPath) {
			return false
		}
	}
	return true
}

func isSubPath(parent, child string) bool {
	parent = strings.TrimSuffix(parent, "/") + "/"
	child = strings.TrimSuffix(child, "/") + "/"
	return strings.HasPrefix(child, parent)
}
