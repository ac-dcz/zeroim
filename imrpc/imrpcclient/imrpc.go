// Code generated by goctl. DO NOT EDIT.
// Source: imrpc.proto

package imrpcclient

import (
	"context"

	"zeroim/imrpc/imrpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Request  = imrpc.Request
	Response = imrpc.Response

	Imrpc interface {
		Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	}

	defaultImrpc struct {
		cli zrpc.Client
	}
)

func NewImrpc(cli zrpc.Client) Imrpc {
	return &defaultImrpc{
		cli: cli,
	}
}

func (m *defaultImrpc) Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	client := imrpc.NewImrpcClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}
