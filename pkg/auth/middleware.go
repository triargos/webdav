package auth

import (
	"context"
	"fmt"
	"github.com/triargos/webdav/pkg/helper"
	"log/slog"
	"net/http"
)

func BasicAuthMiddleware(authenticationService Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			username, password, ok := request.BasicAuth()
			if !ok {
				writer.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
				slog.Error("Unauthorized access attempt: No credentials provided", "remote_addr", request.RemoteAddr)
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if !authenticationService.Authenticate(username, password) {
				slog.Error("Unauthorized access attempt: Invalid credentials", "remote_addr", request.RemoteAddr, "username", username)
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if !authenticationService.HasPermission(request.URL.Path, username) {
				slog.Error("Forbidden access attempt", "remote_addr", request.RemoteAddr, "username", username, "path", request.URL.Path)
				http.Error(writer, "Forbidden", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(request.Context(), helper.UserNameContextKey, username)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

const realm = "WebDAV"

func DigestAuthMiddleware(digestAuthenticator DigestAuthenticator, authenticationService Service) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authHeader := request.Header.Get("Authorization")
			if authHeader == "" {
				slog.Error("Unauthorized access attempt: No credentials provided", "remote_addr", request.RemoteAddr)
				nonce := digestAuthenticator.GenerateNonce()
				writer.Header().Set("WWW-Authenticate", fmt.Sprintf(`Digest realm="%s", qop="auth", nonce="%s", opaque="%s"`, realm, nonce, digestAuthenticator.GenerateOpaque(nonce)))
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}
			username, ok := digestAuthenticator.Authenticate(AuthenticateDigestOptions{
				AuthHeader: authHeader,
				Method:     request.Method,
				Uri:        request.URL.Path,
			})
			if !ok {
				slog.Error("Unauthorized access attempt: Invalid credentials", "remote_addr", request.RemoteAddr, "username", username)
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if !authenticationService.HasPermission(request.URL.Path, username) {
				slog.Error("Forbidden access attempt", "remote_addr", request.RemoteAddr, "username", username, "path", request.URL.Path)
				http.Error(writer, "Forbidden", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(request.Context(), helper.UserNameContextKey, username)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}

}
