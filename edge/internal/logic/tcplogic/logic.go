package tcplogic

import (
	"context"
	"zeroim/core/protocol"
	"zeroim/core/socket"
	"zeroim/edge/internal/svc"

	"github.com/segmentio/kafka-go"
)

type TcpLogic struct {
	svcCtx *svc.ServiceContext
}

func NewTcpLogic(svcCtx *svc.ServiceContext) *TcpLogic {
	return &TcpLogic{
		svcCtx: svcCtx,
	}
}

func (l *TcpLogic) HandleMessage(ctx context.Context, msg *protocol.IMMessage) error {

	return nil
}

func (l *TcpLogic) HandleUserOnLine(ctx context.Context, session *socket.Session) error {

	return nil
}

func (l *TcpLogic) HandleUserOffLine(ctx context.Context, session *socket.Session) error {

	return nil
}

func (l *TcpLogic) HandleKafkaMessage(ctx context.Context, msg *kafka.Message) error {

	return nil
}
