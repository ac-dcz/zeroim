package logic

import (
	"context"
	"zeroim/edge/internal/svc"
	"zeroim/edge/types"
)

type ImrpcLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImrpcLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImrpcLogic {
	return &ImrpcLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImrpcLogic) Connect(msg *types.ConnectMsg) error {

	//Step1. check token

	//Step2. UserID

	//Step3. CreateSession

	//Step4. imrpc call

	return nil
}
