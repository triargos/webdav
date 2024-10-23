package auth

import (
	"context"
	"github.com/triargos/webdav/pkg/cookie"
	"github.com/triargos/webdav/pkg/helper"
	"net/http"
)

type AuthenticationMiddleware struct {
	authenticator Authenticator
	cookieService *cookie.Service
}

func (middleware *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !middleware.shouldAuthenticate(request) {
			next.ServeHTTP(writer, request)
			return
		}
		session, parseSessionErr := middleware.cookieService.ParseSession(request)
		if parseSessionErr != nil {
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
	if request.Method == http.MethodOptions || request.Method == http.MethodHead {
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
