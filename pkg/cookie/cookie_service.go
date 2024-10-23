package cookie

import (
	"fmt"
	"github.com/lucsky/cuid"
	"net/http"
	"time"
)

type Service struct {
	sessions map[string]*Session
}

const name = "webdav_auth"

type Session struct {
	Username  string
	expiresAt time.Time
}

func New() *Service {
	return &Service{
		sessions: make(map[string]*Session),
	}
}

func (service *Service) ParseSession(request *http.Request) (*Session, error) {
	cookie, getCookieErr := request.Cookie(name)
	if getCookieErr != nil {
		return nil, fmt.Errorf("cookie not found")
	}
	sessionToken := cookie.Value
	session, getSessionErr := service.GetSession(sessionToken)
	if getSessionErr != nil {
		return nil, fmt.Errorf("session does not exist")
	}
	if session.IsExpired() {
		service.RemoveSession(sessionToken)
		return nil, fmt.Errorf("session expired")
	}
	return session, nil

}

func (service *Service) CreateSession(username string) (*http.Cookie, *Session) {
	cookieSession := service.generateSession(username)
	sessionToken := service.generateSessionToken()
	service.addSession(cookieSession, sessionToken)
	return &http.Cookie{
		Name:     name,
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  cookieSession.expiresAt,
	}, cookieSession
}

func (service *Service) generateSession(username string) *Session {
	return &Session{
		Username:  username,
		expiresAt: time.Now().Add(365 * 24 * time.Hour),
	}
}

func (service *Service) addSession(session *Session, token string) {
	service.sessions[token] = session
}

func (service *Service) GetSession(token string) (*Session, error) {
	session, ok := service.sessions[token]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (service *Service) RemoveSession(token string) {
	delete(service.sessions, token)
}

func (session Session) IsExpired() bool {
	return session.expiresAt.Before(time.Now())
}

func (service *Service) generateSessionToken() string {
	return cuid.New()
}
