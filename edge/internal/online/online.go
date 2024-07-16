package online

import (
	"sync"
	"zeroim/edge/internal/session"
)

type Manager struct {
	sMutex   sync.RWMutex
	sessions map[int64]*session.Session
}

func NewManager() *Manager {
	return &Manager{
		sMutex:   sync.RWMutex{},
		sessions: make(map[int64]*session.Session),
	}
}

func (m *Manager) AddSession(s *session.Session) {
	m.sMutex.Lock()
	defer m.sMutex.Unlock()
	m.sessions[s.UserID] = s
}

func (m *Manager) DelSession(userId int64) {
	m.sMutex.Lock()
	defer m.sMutex.Unlock()
	delete(m.sessions, userId)
}

func (m *Manager) Session(userId int64) *session.Session {
	m.sMutex.RLock()
	defer m.sMutex.RUnlock()
	return m.sessions[userId]
}
