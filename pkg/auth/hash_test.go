package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triargos/webdav/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

// Test isHashed function
func TestIsHashed(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"$2a$10$VbY8k0P8K/BnWpXc2J2/1uThKnRtrP1KqXe7zxQ2tfq6FbDHE7c5C", true},
		{"$2b$10$VbY8k0P8K/BnWpXc2J2/1uThKnRtrP1KqXe7zxQ2tfq6FbDHE7c5C", true},
		{"$2y$10$VbY8k0P8K/BnWpXc2J2/1uThKnRtrP1KqXe7zxQ2tfq6FbDHE7c5C", true},
		{"plainPassword", false},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			result := isHashed(tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test GenHash function
func TestGenHash(t *testing.T) {
	tests := []struct {
		password string
	}{
		{"testPassword1"},
		{"anotherPassword"},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			hash := GenHash([]byte(tt.password))
			err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
			assert.NoError(t, err, "Hashed password should match the original password")
		})
	}
}

// Test HashPasswords function
func TestHashPasswords(t *testing.T) {
	setupMockUsers(map[string]config.User{
		"testUser": {
			Password: "testPassword",
		}})
	HashPasswords()

	for username, user := range *config.Get().Users {
		t.Run(username, func(t *testing.T) {
			assert.True(t, isHashed(user.Password), "Password should be hashed")
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("testPassword"))
			assert.NoError(t, err, "Hashed password should match the original password")
		})
	}
}
