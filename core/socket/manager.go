package socket

import (
	"fmt"
	"sync"
	"zeroim/core/session"
)

type SessionHandle func(s *Session)

type Manager struct {
	sMutex        sync.Mutex
	sessions      map[session.SessionID]*Session
	BeforeAddFunc SessionHandle
	AfterRemFunc  SessionHandle
	closeFlag     bool
	cMutex        sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[session.SessionID]*Session),
	}
}

func (m *Manager) AddSession(s *Session) error {
	m.sMutex.Lock()
	defer m.sMutex.Unlock()
	if m.BeforeAddFunc != nil {
		go m.BeforeAddFunc(s)
	}
	m.sessions[s.sessionID] = s
	return nil
}

func (m *Manager) RemSession(sessionID session.SessionID) error {
	m.sMutex.Lock()
	defer m.sMutex.Unlock()
	if session, ok := m.sessions[sessionID]; ok {
		if m.AfterRemFunc != nil {
			go m.AfterRemFunc(session)
			delete(m.sessions, sessionID)
		}
	}
	return nil
}

func (m *Manager) GetSession(sessionID session.SessionID) (*Session, bool) {
	m.sMutex.Lock()
	defer m.sMutex.Unlock()
	session, ok := m.sessions[sessionID]
	return session, ok
}

func (m *Manager) Close() error {
	m.cMutex.Lock()
	defer m.cMutex.Unlock()
	if m.closeFlag {
		return fmt.Errorf("manager has been closed")
	}
	for sid, s := range m.sessions {
		m.RemSession(sid)
		s.Close()
	}
	m.closeFlag = true
	return nil
}
