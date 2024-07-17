package auth

import (
	"fmt"
	"github.com/triargos/webdav/mocks"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/helper"
	"strings"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	mockUserService := mocks.NewMockUserService(map[string]config.User{
		"testuser": {Password: helper.Md5Hash("testuser:testpassword:WebDAV")},
	})
	authenticator := NewDigestAuthenticator(mockUserService)

	tests := []struct {
		name     string
		options  AuthenticateDigestOptions
		expected string
		ok       bool
	}{
		{
			name: "Valid authentication",
			options: AuthenticateDigestOptions{
				AuthHeader: `Digest username="testuser", realm="WebDAV", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", uri="/dir/index.html", qop=auth, nc=00000001, cnonce="0a4f113b", response="%s", opaque="5ccc069c403ebaf9f0171e9517f40e41"`,
				Method:     "GET",
				Uri:        "/dir/index.html",
			},
			expected: "testuser",
			ok:       true,
		},
		{
			name: "Invalid realm",
			options: AuthenticateDigestOptions{
				AuthHeader: `Digest username="testuser", realm="InvalidRealm", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", uri="/dir/index.html", qop=auth, nc=00000001, cnonce="0a4f113b", response="6629fae49393a05397450978507c4ef1", opaque="5ccc069c403ebaf9f0171e9517f40e41"`,
				Method:     "GET",
				Uri:        "/dir/index.html",
			},
			expected: "",
			ok:       false,
		},
		{
			name: "User does not exist",
			options: AuthenticateDigestOptions{
				AuthHeader: `Digest username="nonexistent", realm="WebDAV", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", uri="/dir/index.html", qop=auth, nc=00000001, cnonce="0a4f113b", response="6629fae49393a05397450978507c4ef1", opaque="5ccc069c403ebaf9f0171e9517f40e41"`,
				Method:     "GET",
				Uri:        "/dir/index.html",
			},
			expected: "nonexistent",
			ok:       false,
		},
		{
			name: "Invalid response",
			options: AuthenticateDigestOptions{
				AuthHeader: `Digest username="testuser", realm="WebDAV", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", uri="/dir/index.html", qop=auth, nc=00000001, cnonce="0a4f113b", response="invalidresponse", opaque="5ccc069c403ebaf9f0171e9517f40e41"`,
				Method:     "GET",
				Uri:        "/dir/index.html",
			},
			expected: "testuser",
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authHeader := tt.options.AuthHeader
			if tt.ok {
				ha1 := mockUserService.GetUser(tt.expected).Password
				ha2 := helper.Md5Hash(fmt.Sprintf("%s:%s", tt.options.Method, tt.options.Uri))
				params := ParseAuthHeader(authHeader)
				expectedResponse := helper.Md5Hash(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, params["nonce"], params["nc"], params["cnonce"], params["qop"], ha2))
				authHeader = fmt.Sprintf(tt.options.AuthHeader, expectedResponse)
				tt.options.AuthHeader = authHeader
			}

			username, ok := authenticator.Authenticate(tt.options)
			if username != tt.expected || ok != tt.ok {
				t.Errorf("expected %s, %v; got %s, %v", tt.expected, tt.ok, username, ok)
			}
		})
	}
}

func TestGenerateNonce(t *testing.T) {
	authenticator := NewDigestAuthenticator(nil)
	nonce := authenticator.GenerateNonce()
	if nonce == "" {
		t.Error("expected non-empty nonce")
	}
}

func TestGenerateOpaque(t *testing.T) {
	authenticator := NewDigestAuthenticator(nil)
	nonce := authenticator.GenerateNonce()
	opaque := authenticator.GenerateOpaque(nonce)
	if opaque == "" {
		t.Error("expected non-empty opaque")
	}
}

func TestParseAuthHeader(t *testing.T) {
	header := `Digest username="testuser", realm="WebDAV", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", uri="/dir/index.html", qop=auth, nc=00000001, cnonce="0a4f113b", response="6629fae49393a05397450978507c4ef1", opaque="5ccc069c403ebaf9f0171e9517f40e41"`
	expected := map[string]string{
		"username": "testuser",
		"realm":    "WebDAV",
		"nonce":    "dcd98b7102dd2f0e8b11d0f600bfb0c093",
		"uri":      "/dir/index.html",
		"qop":      "auth",
		"nc":       "00000001",
		"cnonce":   "0a4f113b",
		"response": "6629fae49393a05397450978507c4ef1",
		"opaque":   "5ccc069c403ebaf9f0171e9517f40e41",
	}
	params := ParseAuthHeader(strings.TrimPrefix(header, "Digest "))
	for k, v := range expected {
		if params[k] != v {
			t.Errorf("expected %s=%s, got %s", k, v, params[k])
		}
	}
}
