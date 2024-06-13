package auth

import (
	"github.com/triargos/webdav/pkg/logging"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triargos/webdav/pkg/config"
)

func setupMockUsers(users map[string]config.User) {
	config.Set(&config.Config{
		Users: &users,
	})
}

// Test isSubPath function
func TestIsSubPath(t *testing.T) {
	tests := []struct {
		parent   string
		child    string
		expected bool
	}{
		{"/users/admin", "/users/admin/exports", true},
		{"/users/user", "/users/user/exports/file.txt", true},
		{"/users", "/users/admin/", true},
		{"/users/admin", "/users/user", false},
		{"/users/admin", "/users/user/exports", false},
		{"/users/admin/exports", "/users/user/exports/file.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.parent+"-"+tt.child, func(t *testing.T) {
			result := isSubPath(tt.parent, tt.child)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test HasPermission function
func TestHasPermission(t *testing.T) {
	tests := []struct {
		name        string
		users       map[string]config.User
		currentPath string
		username    string
		expected    bool
	}{
		{
			name: "Admin user",
			users: map[string]config.User{
				"user1": {Admin: true},
			},
			currentPath: "/any/path",
			username:    "user1",
			expected:    true,
		},
		{
			name: "Jailed user with permission",
			users: map[string]config.User{
				"user1": {Jail: true, Root: "/home/user1"},
			},
			currentPath: "/home/user1/documents",
			username:    "user1",
			expected:    true,
		},
		{
			name: "Jailed user without permission",
			users: map[string]config.User{
				"user1": {Jail: true, Root: "/home/user1"},
			},
			currentPath: "/home/user2/documents",
			username:    "user1",
			expected:    false,
		},
		{
			name: "Jailed user to root path",
			users: map[string]config.User{
				"user1": {Jail: true, Root: "/home/user1"},
			},
			currentPath: "/home/documents",
			username:    "user1",
			expected:    false,
		},
		{
			name: "User accessing other's root",
			users: map[string]config.User{
				"user1": {Root: "/home/user1"},
				"user2": {Root: "/home/user2"},
			},
			currentPath: "/home/user2/documents",
			username:    "user1",
			expected:    false,
		},
		{
			name: "User accessing own root",
			users: map[string]config.User{
				"user1": {Root: "/home/user1"},
			},
			currentPath: "/home/user1/documents",
			username:    "user1",
			expected:    true,
		},
		{
			name: "User with no specific permissions",
			users: map[string]config.User{
				"user1": {},
			},
			currentPath: "/any/path",
			username:    "user1",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logging.InitLoggers()
			setupMockUsers(tt.users)
			result := HasPermission(tt.currentPath, tt.username)
			assert.Equal(t, tt.expected, result)
		})
	}
}
