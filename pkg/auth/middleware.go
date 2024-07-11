package auth

import (
	"context"
	"github.com/triargos/webdav/pkg/helper"
	"log/slog"
	"net/http"
)

func Middleware(authenticationService Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			username, password, ok := request.BasicAuth()
			ctx := context.WithValue(request.Context(), helper.UserNameContextKey, username)
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
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
