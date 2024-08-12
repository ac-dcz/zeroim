package etcd

import (
	"zeroim/common/sd"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Discover struct {
	cli *client
}

func NewDiscover(cfg clientv3.Config) (sd.Discover, error) {
	cli, err := newClient(cfg, defaultTTL)
	if err != nil {
		return nil, err
	}
	return &Discover{
		cli: cli,
	}, nil
}

func (d *Discover) Endpoints(name string) ([]sd.Service, error) {
	return d.cli.endpoints(name)
}
