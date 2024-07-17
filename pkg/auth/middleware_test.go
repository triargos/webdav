package auth_test

import (
	"fmt"
	"github.com/triargos/webdav/mocks"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/helper"
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
			urlPath:        "/Finanzen/mydir",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Non-admin user with valid credentials and permission",
			users: map[string]config.User{
				"user1": {Password: string(hash), Admin: false},
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
			req, err := http.NewRequest("GET", tt.urlPath, nil)
			assert.NoError(t, err)
			if tt.username != "" && tt.password != "" {
				req.SetBasicAuth(tt.username, tt.password)
			}
			rr := httptest.NewRecorder()
			userService := mocks.NewMockUserService(tt.users)
			authenticationService := auth.New(userService)
			handler := auth.BasicAuthMiddleware(authenticationService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestDigestAuthMiddleware(t *testing.T) {

	mockUserService := mocks.NewMockUserService(map[string]config.User{
		"testuser": {Password: helper.Md5Hash("testuser:WebDAV:testpassword"), Jail: true, Root: "/users/testuser"},
	})

	mockAuthService := auth.New(mockUserService)

	digestAuthenticator := auth.NewDigestAuthenticator(mockUserService)
	middleware := auth.DigestAuthMiddleware(digestAuthenticator, mockAuthService)

	tests := []struct {
		name          string
		authHeader    string
		requestPath   string
		expectedCode  int
		expectedLog   string
		expectedNonce string
	}{
		{
			name:         "No Auth Header",
			authHeader:   "",
			requestPath:  "/allowed",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "Unauthorized access attempt: No credentials provided",
		},
		{
			name:         "Invalid Credentials",
			authHeader:   `Digest username="testuser", realm="WebDAV", nonce="invalidnonce", uri="/allowed", qop=auth, nc=00000001, cnonce="0a4f113b", response="invalidresponse", opaque="5ccc069c403ebaf9f0171e9517f40e41"`,
			requestPath:  "/allowed",
			expectedCode: http.StatusUnauthorized,
			expectedLog:  "Unauthorized access attempt: Invalid credentials",
		},
		{
			name:         "Valid Credentials, No Permission",
			authHeader:   `Digest username="testuser", realm="WebDAV", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", uri="/forbidden", qop=auth, nc=00000001, cnonce="0a4f113b", response="%s", opaque="5ccc069c403ebaf9f0171e9517f40e41"`,
			requestPath:  "/forbidden",
			expectedCode: http.StatusForbidden,
			expectedLog:  "Forbidden access attempt",
		},
		{
			name:         "Valid Credentials, Has Permission",
			authHeader:   `Digest username="testuser", realm="WebDAV", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", uri="/allowed", qop=auth, nc=00000001, cnonce="0a4f113b", response="%s", opaque="5ccc069c403ebaf9f0171e9517f40e41"`,
			requestPath:  "/users/testuser/mydir",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request and response
			req, err := http.NewRequest("GET", tt.requestPath, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			if tt.name == "Valid Credentials, No Permission" || tt.name == "Valid Credentials, Has Permission" {
				ha1 := mockUserService.GetUser("testuser").Password
				ha2 := helper.Md5Hash(fmt.Sprintf("GET:%s", tt.requestPath))
				params := auth.ParseAuthHeader(tt.authHeader)
				expectedResponse := helper.Md5Hash(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, params["nonce"], params["nc"], params["cnonce"], params["qop"], ha2))
				tt.authHeader = fmt.Sprintf(tt.authHeader, expectedResponse)
				req.Header.Set("Authorization", tt.authHeader)
			}

			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}
		})
	}
}
