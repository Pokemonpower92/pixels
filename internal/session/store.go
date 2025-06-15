package session

import (
	"sync"

	"github.com/google/uuid"
)

type Sessionizer interface {
	CreateSession(uuid.UUID) uuid.UUID
	FindSession(uuid.UUID) (bool, uuid.UUID)
}

type SessionStore struct {
	sessions map[string]string
	mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]string),
	}
}

func (ss *SessionStore) CreateSession(user_id uuid.UUID) uuid.UUID {
	sessionId := uuid.New()
	ss.mu.Lock()
	ss.sessions[sessionId.String()] = user_id.String()
	ss.mu.Unlock()
	return sessionId
}

func (ss *SessionStore) FindSession(sessionId uuid.UUID) (bool, uuid.UUID) {
	ss.mu.RLock()
	userId, exists := ss.sessions[sessionId.String()]
	ss.mu.RUnlock()
	if !exists {
		return false, uuid.UUID{}
	}
	return true, uuid.MustParse(userId)
}
