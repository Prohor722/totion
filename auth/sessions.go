package auth

import (
	"sync"

	"github.com/Prohor722/totion/model"
)

// SessionRepository defines operations for session persistence used by auth.
type SessionRepository interface {
	Create(session *model.Session)
	Get(sessionID string) (*model.Session, bool)
	Remove(sessionID string)
}

// InMemorySessionRepository is a simple in-memory implementation for sessions.
type InMemorySessionRepository struct {
	mutex    sync.RWMutex
	sessions map[string]*model.Session
}

// NewInMemorySessionRepository constructs an in-memory session store.
func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]*model.Session)}
}

func (r *InMemorySessionRepository) Create(session *model.Session) {
	if session == nil {
		return
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sessions[session.SessionID] = session
}

func (r *InMemorySessionRepository) Get(sessionID string) (*model.Session, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	session, ok := r.sessions[sessionID]
	return session, ok
}

func (r *InMemorySessionRepository) Remove(sessionID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.sessions, sessionID)
}
