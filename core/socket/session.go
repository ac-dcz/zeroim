package socket

import (
	"encoding/json"
	"fmt"
	"sync"
	"zeroim/core/jwt"
	"zeroim/core/protocol"
	"zeroim/core/session"
)

type Session struct {
	EdgeId    uint64
	Uid       uint64
	cc        protocol.Codec
	manager   *Manager
	sessionID session.SessionID
	tokenOpt  *jwt.TokenOption
	closeFlag bool
	cMutex    sync.Mutex
}

func NewSession(edgeId uint64, cc protocol.Codec, manager *Manager, tokenOpt *jwt.TokenOption) *Session {
	return &Session{
		cc:       cc,
		manager:  manager,
		tokenOpt: tokenOpt,
		EdgeId:   edgeId,
	}
}

func (s *Session) run() error {
	//Step1: Session Hand
	data, err := s.cc.ShakeHand(s.tokenOpt)
	if err != nil {
		return err
	}
	uid, err := data["uid"].(json.Number).Int64()
	if err != nil {
		return err
	}
	s.Uid = uint64(uid)
	s.sessionID = session.SessionID{
		Uid:    s.Uid,
		EdgeId: s.EdgeId,
	}

	//Step2: Add manager

	//Step3: Message Loop

	return nil
}

func (s *Session) Close() error {
	s.cMutex.Lock()
	defer s.cMutex.Unlock()
	if s.closeFlag {
		return fmt.Errorf("session has been closed")
	}
	s.closeFlag = true
	return s.cc.Close()
}
