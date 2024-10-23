package auth

import (
	"context"
	"github.com/triargos/webdav/pkg/cookie"
	"github.com/triargos/webdav/pkg/helper"
	"log/slog"
	"net/http"
)

type AuthenticationMiddleware struct {
	authenticator Authenticator
	cookieService *cookie.Service
}

func (middleware *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.shouldAuthenticate(request) {
			slog.Info("Skipping authentication for request", "method", request.Method, "path", request.URL.Path)
			next.ServeHTTP(writer, request)
			return
		}
		slog.Info("Performing authentication for request", "method", request.Method, "path", request.URL.Path)
		session, parseSessionErr := middleware.cookieService.ParseSession(request)
		if parseSessionErr != nil {
			slog.Info("Error parsing session", "error", parseSessionErr)
			authenticatedUserName, updatedResponseWriter := middleware.authenticator.PerformAuthentication(writer, request)
			if authenticatedUserName == "" {
				http.Error(updatedResponseWriter, "Unauthorized", http.StatusUnauthorized)
				return
			}
			sessionCookie, createdSession := middleware.cookieService.CreateSession(authenticatedUserName)
			session = createdSession
			http.SetCookie(writer, sessionCookie)
		}
		ctx := context.WithValue(request.Context(), helper.UserNameContextKey, session.Username)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (middleware *AuthenticationMiddleware) shouldAuthenticate(request *http.Request) bool {
	if request.Method == http.MethodOptions || request.Method == http.MethodGet {
		return false
	}
	return true
}

func NewMiddleware(authenticator Authenticator, cookieService *cookie.Service) AuthenticationMiddleware {
	return AuthenticationMiddleware{
		authenticator: authenticator,
		cookieService: cookieService,
	}
}
