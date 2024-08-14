package socket

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"zeroim/core/jwt"
	"zeroim/core/protocol"
	"zeroim/core/session"
)

type Session struct {
	EdgeId    uint64
	Uid       uint64
	msgChan   chan *protocol.IMMessage
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
		msgChan:  make(chan *protocol.IMMessage, 1000),
	}
}

func (s *Session) SendMessage(msg *protocol.IMMessage) {
	s.msgChan <- msg
}

func (s *Session) SessionID() session.SessionID {
	return s.sessionID
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
	s.manager.AddSession(s)
	//Step3: Message Loop
	go s.messageLoop()

	go s.healthListen()

	return nil
}

func (s *Session) messageLoop() {
	defer s.Close()
	for msg := range s.msgChan {
		if err := s.cc.Write(msg); err != nil {
			log.Printf("Write Message Error: %v\n", err)
			return
		}
	}
}

func (s *Session) healthListen() {
	defer s.Close()
	for {
		_, err := s.cc.Receive()
		if err == io.EOF {
			fmt.Printf("Connection %d is closed\n", s.Uid)
			return
		} else if err != nil {
			fmt.Printf("Connection %d error: %v\n", s.Uid, err)
			return
		}
	}
}

func (s *Session) Close() error {
	s.cMutex.Lock()
	defer s.cMutex.Unlock()
	if s.closeFlag {
		return fmt.Errorf("session has been closed")
	}
	s.closeFlag = true

	//RemSession
	s.manager.RemSession(s.sessionID)
	return s.cc.Close()
}
