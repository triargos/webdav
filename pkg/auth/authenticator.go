package auth

import (
	"context"
	"fmt"
	"github.com/triargos/webdav/pkg/helper"
	"github.com/triargos/webdav/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
)

type Authenticator interface {
	Middleware(next http.Handler) http.Handler
}

type BasicAuthenticator struct {
	userService user.Service
}

func (authenticator *BasicAuthenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if authenticator.isAllowedMethod(request) {
			next.ServeHTTP(writer, request)
			return
		}
		validateCookieErr, isPresent := authenticator.validateCookie(request)
		if validateCookieErr == nil {
			slog.Info("Validated cookie", "path", request.URL.Path, "method", request.Method, "remote_addr", request.RemoteAddr)
			next.ServeHTTP(writer, request)
			return
		} else if isPresent {
			slog.Error("Failed to validate cookie", "error", validateCookieErr)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		slog.Info("No cookie found", "path", request.URL.Path, "method", request.Method, "remote_addr", request.RemoteAddr)

		//Otherwise, we did not authenticate yet
		username, password, ok := request.BasicAuth()
		if !ok {
			writer.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
			slog.Info("Requesting basic auth", "path", request.URL.Path, "method", request.Method, "remote_addr", request.RemoteAddr)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		validateCredentialsErr := authenticator.validateCredentials(username, password)
		if validateCredentialsErr != nil {
			slog.Error("Failed to validate credentials", "error", validateCredentialsErr)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		slog.Info("Setting session cookie", "path", request.URL.Path, "method", request.Method, "remote_addr", request.RemoteAddr)
		authenticator.addSessionCookie(username, writer)
		//TODO: Validate permissions for user
		userCtx := context.WithValue(request.Context(), helper.UserNameContextKey, username)
		next.ServeHTTP(writer, request.WithContext(userCtx))
	})
}

func (authenticator *BasicAuthenticator) validateCredentials(username, password string) error {
	if !authenticator.userService.HasUser(username) {
		return fmt.Errorf("user not found")
	}
	userObject := authenticator.userService.GetUser(username)
	verifyPasswordErr := bcrypt.CompareHashAndPassword([]byte(userObject.Password), []byte(password))
	return verifyPasswordErr
}

func (authenticator *BasicAuthenticator) addSessionCookie(username string, writer http.ResponseWriter) {
	session := NewSession(username)
	sessionToken := generateSessionToken()
	GlobalSessionManager.AddSession(sessionToken, session)
	cookie := http.Cookie{
		Name:     "webdav_auth",
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  session.expiresAt,
	}
	http.SetCookie(writer, &cookie)

}

func (authenticator *BasicAuthenticator) validateCookie(request *http.Request) (error, bool) {
	cookie, getCookieErr := request.Cookie("webdav_auth")
	if getCookieErr != nil {
		return fmt.Errorf("no cookie found"), false
	}
	sessionToken := cookie.Value
	userSession, getSessionErr := GlobalSessionManager.GetSession(sessionToken)
	if getSessionErr != nil {
		return fmt.Errorf("failed to get session: %v", getSessionErr), false
	}
	if userSession.IsExpired() {
		GlobalSessionManager.RemoveSession(sessionToken)
		return fmt.Errorf("session expired"), true
	}
	userCtx := context.WithValue(request.Context(), helper.UserNameContextKey, userSession.Username)
	request = request.WithContext(userCtx)
	return nil, true
}

func (authenticator *BasicAuthenticator) isAllowedMethod(request *http.Request) bool {
	return request.Method == http.MethodHead || request.Method == http.MethodOptions
}

func NewBasicAuthenticator(userService user.Service) Authenticator {
	return &BasicAuthenticator{userService: userService}
}
