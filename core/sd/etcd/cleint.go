package etcd

import (
	"context"
	"errors"
	"zeroim/common/sd"

	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

var (
	ErrNameEmpty = errors.New("server name is empty")
	ErrAddrEmpty = errors.New("server address is empty")
)

type client struct {
	cli    *clientv3.Client
	ttlOpt TTLOption
}

func newClient(cfg clientv3.Config, ttl TTLOption) (*client, error) {
	c := &client{ttlOpt: ttl}
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	c.cli = cli
	return c, nil
}

func (c *client) registry(s *sd.Service) error {
	if s.Name == "" {
		return ErrNameEmpty
	}
	if s.Addr == "" {
		return ErrAddrEmpty
	}

	resp, err := c.cli.Lease.Grant(context.Background(), int64(c.ttlOpt.TTL.Seconds()))
	if err != nil {
		return err
	}
	manager, err := endpoints.NewManager(c.cli, s.Name)
	if err != nil {
		return err
	}
	val, _ := s.Encode()
	if err := manager.AddEndpoint(context.Background(), s.Key(), endpoints.Endpoint{
		Addr:     s.Addr,
		Metadata: val,
	}, clientv3.WithLease(resp.ID)); err != nil {
		return err
	}

	lch, err := c.cli.Lease.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	go func() {
		for q := range lch {
			logx.Debugf("%s Lease keep live ID: %d", s.Key(), q.ID)
		}
		logx.Debugf("%s Lease is expired", s.Key())
	}()

	return nil
}

func (c *client) disregistry(s *sd.Service) error {
	if s.Name == "" {
		return ErrNameEmpty
	}
	if s.Addr == "" {
		return ErrAddrEmpty
	}
	manager, err := endpoints.NewManager(c.cli, s.Name)
	if err != nil {
		return err
	}
	return manager.DeleteEndpoint(context.Background(), s.Key(), clientv3.WithIgnoreLease())
}

func (c *client) endpoints(name string) ([]sd.Service, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}
	manager, err := endpoints.NewManager(c.cli, name)
	if err != nil {
		return nil, err
	}
	ends, err := manager.List(context.Background())
	if err != nil {
		return nil, err
	}
	var ret []sd.Service
	for _, endponit := range ends {
		s := &sd.Service{}
		if err := s.Decode(endponit.Metadata.([]byte)); err != nil {
			return nil, err
		}
		ret = append(ret, *s)
	}
	return ret, nil
}
