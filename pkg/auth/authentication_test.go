package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triargos/webdav/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

// Test AuthenticateUser function
func TestAuthenticateUser(t *testing.T) {
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
			setupMockUsers(tt.users)
			result := AuthenticateUser(tt.username, tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}
