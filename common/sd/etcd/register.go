package etcd

import (
	"time"
	"zeroim/common/sd"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type TTLOption struct {
	TTL       time.Duration //Lease: time to live
	HeartBeat time.Duration //用于检测服务是否还存活
}

var defaultTTL = TTLOption{
	TTL:       10 * time.Second, // 10s
	HeartBeat: 5 * time.Second,  //5s
}

type Option func(r *RegisterOption)

func WithTTL(ttl TTLOption) Option {
	return func(r *RegisterOption) {
		r.ttl = ttl
	}
}

type RegisterOption struct {
	ttl TTLOption
}

type Register struct {
	cli    *client
	regOpt *RegisterOption
}

func NewRegister(cfg clientv3.Config, opts ...Option) (sd.Register, error) {
	r := &Register{
		regOpt: &RegisterOption{},
	}
	for _, opt := range opts {
		opt(r.regOpt)
	}
	cli, err := newClient(cfg, r.regOpt.ttl)
	if err != nil {
		return nil, err
	}
	r.cli = cli
	return r, nil
}

func (r *Register) Registry(s *sd.Service) error {
	return r.cli.registry(s)
}

func (r *Register) DisRegistry(s *sd.Service) error {
	return r.cli.disregistry(s)
}
