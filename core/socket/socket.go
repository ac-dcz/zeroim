package socket

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
	"zeroim/core/jwt"
	"zeroim/core/protocol"
)

type Server struct {
	EdgeId    uint64
	NetWork   string
	Address   string
	protocol  protocol.Protocol
	manager   *Manager
	tlsConfig *tls.Config
	listener  net.Listener
	tokenOpt  *jwt.TokenOption
	close     bool
	cMutex    sync.Mutex
}

func NewServer(manager *Manager, edgeId uint64, p protocol.Protocol) *Server {
	if p == nil {
		p = protocol.ImProtocol{}
	}
	return &Server{
		EdgeId:   edgeId,
		manager:  manager,
		protocol: p,
		close:    false,
	}
}

func (s *Server) Serve(listener net.Listener) error {
	s.listener = listener
	s.NetWork = listener.Addr().Network()
	s.Address = listener.Addr().String()
	return s.run()
}

func (s *Server) ListenAndServe(network, address string, tlsCfg *tls.Config) (err error) {
	if network != "tcp" {
		return fmt.Errorf("NetWork Invaild")
	}
	var listener net.Listener
	if tlsCfg != nil {
		listener, err = tls.Listen(network, address, tlsCfg)
		if err != nil {
			return err
		}
	} else {
		listener, err = net.Listen(network, address)
		if err != nil {
			return err
		}
	}
	s.tlsConfig = tlsCfg
	return s.Serve(listener)
}

func (s *Server) run() error {
	var tempDelay time.Duration
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ne, err := err.(net.Error); err && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return io.EOF
			}
			return err
		}
		session := NewSession(s.EdgeId, s.protocol.NewCodec(conn), s.manager, s.tokenOpt)
		go session.run()
	}
}

func (s *Server) Close() error {
	s.cMutex.Lock()
	defer s.cMutex.Unlock()
	if s.close {
		return fmt.Errorf("Server has been closed")
	}
	s.close = true
	return s.listener.Close()
}
