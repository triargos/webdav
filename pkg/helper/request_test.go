package helper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsernameFromContext(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		expected   string
		shouldFind bool
	}{
		{
			name:       "User found in context",
			ctx:        context.WithValue(context.Background(), UserNameContextKey, "testuser"),
			expected:   "testuser",
			shouldFind: true,
		},
		{
			name:       "User not found in context",
			ctx:        context.Background(),
			expected:   "",
			shouldFind: false,
		},
		{
			name:       "User found but not a string",
			ctx:        context.WithValue(context.Background(), UserNameContextKey, 12345),
			expected:   "",
			shouldFind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, found := GetUsernameFromContext(tt.ctx)
			assert.Equal(t, tt.shouldFind, found)
			assert.Equal(t, tt.expected, username)
		})
	}
}
