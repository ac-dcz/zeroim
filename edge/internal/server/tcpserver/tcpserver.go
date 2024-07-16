package tcpserver

import "zeroim/edge/internal/svc"

type Server struct {
	svcCtx *svc.ServiceContext
}

func NewServer(svcCtx *svc.ServiceContext) *Server {
	return &Server{
		svcCtx: svcCtx,
	}
}

func (s *Server) Serve() error {

	return nil
}
