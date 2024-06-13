package auth

import (
	"context"
	"github.com/triargos/webdav/pkg/helper"
	"github.com/triargos/webdav/pkg/logging"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		ctx := context.WithValue(request.Context(), helper.UserNameContextKey, username)
		if !ok {
			writer.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
			logging.Log.Error.Printf("Unauthorized access attempt from %s\n: No credentials", request.RemoteAddr)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !AuthenticateUser(username, password) {
			logging.Log.Error.Printf("Unauthorized access attempt from %s\n: Invalid credentials", request.RemoteAddr)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !HasPermission(request.URL.Path, username) {
			logging.Log.Error.Printf("Forbidden access attempt from %s\n", request.RemoteAddr)
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
