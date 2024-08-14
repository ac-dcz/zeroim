package svc

import (
	"zeroim/core/socket"
	"zeroim/edge/internal/config"
	"zeroim/edge/internal/mq"
)

type ServiceContext struct {
	Conf    *config.Config
	Manager *socket.Manager
	Kq      *mq.Reader
}

func NewServiceContext(conf *config.Config) (*ServiceContext, error) {
	svcCtx := &ServiceContext{
		Conf:    conf,
		Manager: socket.NewManager(),
		Kq:      mq.NewReader(conf.Kq.Brokers, conf.Kq.Topic, conf.Kq.GroupID),
	}
	return svcCtx, nil
}
