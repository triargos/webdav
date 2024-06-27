package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triargos/webdav/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

func TestMiddleware(t *testing.T) {
	password := "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		users          map[string]config.User
		username       string
		password       string
		urlPath        string
		expectedStatus int
	}{
		{
			name: "Valid credentials and permission",
			users: map[string]config.User{
				"user1": {Password: string(hash), Admin: true},
			},
			username:       "user1",
			password:       password,
			urlPath:        "/",
			expectedStatus: http.StatusOK,
		},
		{
			name: "No credentials",
			users: map[string]config.User{
				"user1": {Password: string(hash), Admin: true},
			},
			username:       "",
			password:       "",
			urlPath:        "/",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid credentials",
			users: map[string]config.User{
				"user1": {Password: string(hash), Admin: true},
			},
			username:       "user1",
			password:       "wrongPassword",
			urlPath:        "/",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid credentials but no permission",
			users: map[string]config.User{
				"user1": {Password: string(hash), Admin: false, Jail: true, Root: "/home/user1"},
			},
			username:       "user1",
			password:       password,
			urlPath:        "/home/user2",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupMockUsers(tt.users)

			req, err := http.NewRequest("GET", tt.urlPath, nil)
			assert.NoError(t, err)

			if tt.username != "" && tt.password != "" {
				req.SetBasicAuth(tt.username, tt.password)
			}

			rr := httptest.NewRecorder()

			handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
