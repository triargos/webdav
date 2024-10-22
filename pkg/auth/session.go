package auth

import (
	"fmt"
	"github.com/lucsky/cuid"
	"time"
)

var GlobalSessionManager = &SessionManager{
	currentSessions: make(map[string]Session),
}

type SessionManager struct {
	currentSessions map[string]Session
}

func (manager *SessionManager) AddSession(token string, session Session) {
	manager.currentSessions[token] = session
}

func (manager *SessionManager) GetSession(token string) (*Session, error) {
	session, ok := manager.currentSessions[token]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return &session, nil
}

func (manager *SessionManager) RemoveSession(token string) {
	delete(manager.currentSessions, token)
}

type Session struct {
	Username  string
	expiresAt time.Time
}

func NewSession(username string) Session {
	return Session{
		Username:  username,
		expiresAt: time.Now().Add(365 * 24 * time.Hour),
	}
}

func (session Session) IsExpired() bool {
	return session.expiresAt.Before(time.Now())
}

func generateSessionToken() string {
	return cuid.New()
}
