package auth

import (
	"context"
	"fmt"
	"github.com/triargos/webdav/pkg/cookie"
	"github.com/triargos/webdav/pkg/environment"
	"github.com/triargos/webdav/pkg/helper"
	"log/slog"
	"net/http"
	"regexp"
)

type AuthenticationMiddleware struct {
	authenticator Authenticator
	cookieService *cookie.Service
	envService    environment.Service
}

func (middleware *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		middleware.logRequest(request)
		if !middleware.shouldAuthenticate(request) {
			next.ServeHTTP(writer, request)
			return
		}

		session, parseSessionErr := middleware.cookieService.ParseSession(request)
		if parseSessionErr != nil {
			authenticatedUserName, updatedResponseWriter := middleware.authenticator.PerformAuthentication(writer, request)
			if authenticatedUserName == "" {
				http.Error(updatedResponseWriter, "401 Unauthorized", http.StatusUnauthorized)
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

func (middleware *AuthenticationMiddleware) logRequest(request *http.Request) {
	if request.Method == "PROPFIND" {
		return
	}
	slog.Info(fmt.Sprintf("%s %s %s", request.Method, request.URL.Path, request.Header.Get("User-Agent")))
	if middleware.shouldAuthenticate(request) {
		//slog.Warn("AUTH_REQUIRED", "method", request.Method, "path", request.URL.Path)
	}

}

func (middleware *AuthenticationMiddleware) shouldAuthenticate(request *http.Request) bool {
	if request.Method == http.MethodOptions || request.Method == http.MethodHead {
		return false
	}
	regex := regexp.MustCompile("^Microsoft Office(?: .*)?$")
	userAgent := request.Header.Get("User-Agent")
	isOffice := regex.MatchString(userAgent)

	if middleware.envService.GetBool("DISABLE_OFFICE_AUTH") && isOffice {
		return false
	}

	return true
}

func NewMiddleware(authenticator Authenticator, cookieService *cookie.Service, envService environment.Service) AuthenticationMiddleware {
	return AuthenticationMiddleware{
		authenticator: authenticator,
		cookieService: cookieService,
		envService:    envService,
	}
}
