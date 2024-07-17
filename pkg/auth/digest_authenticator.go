package auth

import (
	"fmt"
	"github.com/triargos/webdav/pkg/helper"
	"github.com/triargos/webdav/pkg/user"
	"log/slog"
	"strings"
	"time"
)

type DigestAuthenticator struct {
	userService user.Service
}

type AuthenticateDigestOptions struct {
	AuthHeader string
	Method     string
	Uri        string
}

func NewDigestAuthenticator(userService user.Service) DigestAuthenticator {
	return DigestAuthenticator{userService: userService}
}

func (digestAuthenticator DigestAuthenticator) Authenticate(options AuthenticateDigestOptions) (username string, ok bool) {
	const prefix = "Digest "
	if !strings.HasPrefix(options.AuthHeader, prefix) {
		slog.Error("missing prefix", "prefix", prefix, "auth_header", options.AuthHeader)
		return "", false
	}
	params := ParseAuthHeader(strings.TrimPrefix(options.AuthHeader, prefix))
	if params["realm"] != "WebDAV" {
		slog.Error("invalid realm", "realm", params["realm"])
		return "", false
	}
	username = params["username"]
	if !digestAuthenticator.userService.HasUser(username) {
		slog.Error("user does not exist", "username", username)
		return username, false
	}
	webdavUser := digestAuthenticator.userService.GetUser(username)
	ha1 := webdavUser.Password
	ha2 := helper.Md5Hash(fmt.Sprintf("%s:%s", options.Method, options.Uri))
	response := helper.Md5Hash(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, params["nonce"], params["nc"], params["cnonce"], params["qop"], ha2))
	if response != params["response"] {
		slog.Error("invalid response", "response", response, "expected", params["response"])
		return username, false
	}
	return username, true
}

func (digestAuthenticator DigestAuthenticator) GenerateNonce() string {
	data := fmt.Sprintf("%d", time.Now().UnixNano())
	return helper.Md5Hash(data)
}

func (digestAuthenticator DigestAuthenticator) GenerateOpaque(nonce string) string {
	return helper.Md5Hash(nonce)
}

func ParseAuthHeader(header string) map[string]string {
	params := make(map[string]string)
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			params[strings.Trim(kv[0], `"`)] = strings.Trim(kv[1], `"`)
		}
	}
	return params
}
