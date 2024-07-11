package auth_test

import (
	"github.com/triargos/webdav/mocks"
	"github.com/triargos/webdav/pkg/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triargos/webdav/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

// Test AuthenticateUser function
func TestAuthentication(t *testing.T) {
	password1 := "password123"
	password2 := "anotherPassword"

	hash1, _ := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	hash2, _ := bcrypt.GenerateFromPassword([]byte(password2), bcrypt.DefaultCost)

	tests := []struct {
		name     string
		users    map[string]config.User
		username string
		password string
		expected bool
	}{
		{
			name: "Valid user and password",
			users: map[string]config.User{
				"user1": {Password: string(hash1)},
			},
			username: "user1",
			password: password1,
			expected: true,
		},
		{
			name: "Invalid password",
			users: map[string]config.User{
				"user1": {Password: string(hash1)},
			},
			username: "user1",
			password: "wrongPassword",
			expected: false,
		},
		{
			name: "Non-existing user",
			users: map[string]config.User{
				"user1": {Password: string(hash1)},
			},
			username: "user2",
			password: password1,
			expected: false,
		},
		{
			name: "Valid second user and password",
			users: map[string]config.User{
				"user1": {Password: string(hash1)},
				"user2": {Password: string(hash2)},
			},
			username: "user2",
			password: password2,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := mocks.NewMockUserService(tt.users)
			authenticationService := auth.New(userService)
			result := authenticationService.Authenticate(tt.username, tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test AuthenticateUser function
func TestPermissionCheck(t *testing.T) {

	tests := []struct {
		name     string
		path     string
		username string
		users    map[string]config.User
		expected bool
	}{
		{
			name: "Admin user",
			users: map[string]config.User{
				"user1": {Admin: true},
			},
			username: "user1",
			path:     "/any/path",
			expected: true,
		},
		{
			name: "Valid permissions",
			users: map[string]config.User{
				"user1": {Admin: false, Root: "/Users/user1"},
			},
			username: "user1",
			path:     "/user1/some/dir",
			expected: true,
		},
		{
			name: "Invalid permissions",
			users: map[string]config.User{
				"user2": {Admin: false, Root: "/Users/user2", Jail: true},
			},
			username: "user2",
			path:     "/Users/user1/some/dir",
			expected: false,
		},
		{
			name: "Invalid user",
			users: map[string]config.User{
				"user1": {Admin: false},
			},
			username: "user5",
			path:     "/any/path",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := mocks.NewMockUserService(tt.users)
			authenticationService := auth.New(userService)
			result := authenticationService.HasPermission(tt.path, tt.username)
			assert.Equal(t, tt.expected, result)
		})
	}
}
