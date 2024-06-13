package auth

import (
	"github.com/triargos/webdav/pkg/config"
	"strings"
)

func HasPermission(currentPath, username string) bool {
	cfg := config.Get()
	user, ok := (*cfg.Users)[username]
	if !ok {
		return false
	}
	if user.Admin {
		return true
	}
	if user.Jail {
		return isSubPath(user.Root, currentPath)
	}
	for otherUserName, otherUser := range *cfg.Users {
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
