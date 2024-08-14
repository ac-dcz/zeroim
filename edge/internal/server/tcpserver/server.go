package tcpserver

import (
	"context"
	"zeroim/core/socket"
	"zeroim/edge/internal/logic/tcplogic"
	"zeroim/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/threading"
)

type Server struct {
	svcCtx *svc.ServiceContext
}

func NewServer(svcCtx *svc.ServiceContext) *Server {
	return &Server{
		svcCtx: svcCtx,
	}
}

func (s *Server) RegistryHandle() {
	s.svcCtx.Manager.AddBeforeAddFunc(func(session *socket.Session) {
		l := tcplogic.NewTcpLogic(s.svcCtx)
		l.HandleUserOnLine(context.Background(), session)
	})
	s.svcCtx.Manager.AddAfterRemFunc(func(session *socket.Session) {
		l := tcplogic.NewTcpLogic(s.svcCtx)
		l.HandleUserOffLine(context.Background(), session)
	})

	threading.GoSafe(
		func() {
			l := tcplogic.NewTcpLogic(s.svcCtx)
			ch := s.svcCtx.Kq.MessageChannel()
			for msg := range ch {
				l.HandleKafkaMessage(context.Background(), msg)
			}
		},
	)
}
