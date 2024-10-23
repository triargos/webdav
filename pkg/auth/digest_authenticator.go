package auth

import (
	"fmt"
	"github.com/triargos/webdav/pkg/cookie"
	"github.com/triargos/webdav/pkg/helper"
	"github.com/triargos/webdav/pkg/user"
	"net/http"
	"strings"
	"time"
)

type DigestAuthenticator struct {
	userService   user.Service
	cookieService cookie.Service
}

type AuthenticateDigestOptions struct {
	AuthHeader string
	Method     string
	Uri        string
}

func NewDigestAuthenticator(userService user.Service) DigestAuthenticator {
	return DigestAuthenticator{userService: userService}
}

func (authenticator DigestAuthenticator) PerformAuthentication(writer http.ResponseWriter, request *http.Request) (string, http.ResponseWriter) {
	const realm = "WebDAV"
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		nonce := authenticator.generateNonce()
		writer.Header().Set("WWW-Authenticate", fmt.Sprintf(`Digest realm="%s", qop="auth", nonce="%s", opaque="%s"`, realm, nonce, authenticator.generateOpaque(nonce)))
		return "", writer
	}
	params := authenticator.parseAuthHeader(authHeader)
	username := params["username"]
	if !authenticator.userService.HasUser(username) {
		return "", writer
	}
	validateCredentialsErr := authenticator.validateCredentials(username, request.Method, request.URL.Path, params)
	if validateCredentialsErr != nil {
		return "", writer
	}
	return username, writer

}

func (authenticator DigestAuthenticator) validateCredentials(username, method, uri string, params map[string]string) error {
	webdavUser := authenticator.userService.GetUser(username)
	userPasswordHash := webdavUser.Password
	requestHash := helper.Md5Hash(fmt.Sprintf("%s:%s", method, uri))
	response := helper.Md5Hash(fmt.Sprintf("%s:%s:%s:%s:%s:%s", userPasswordHash, params["nonce"], params["nc"], params["cnonce"], params["qop"], requestHash))
	if response != params["response"] {
		return fmt.Errorf("invalid response")
	}
	return nil
}

func (authenticator DigestAuthenticator) generateNonce() string {
	data := fmt.Sprintf("%d", time.Now().UnixNano())
	return helper.Md5Hash(data)
}

func (authenticator DigestAuthenticator) generateOpaque(nonce string) string {
	return helper.Md5Hash(nonce)
}

func (authenticator DigestAuthenticator) parseAuthHeader(header string) map[string]string {
	params := make(map[string]string)
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			params[strings.Trim(kv[0], `"`)] = strings.Trim(kv[1], `"`)
		}
	}
	return params
}
