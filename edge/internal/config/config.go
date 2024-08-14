package config

import "github.com/zeromicro/go-zero/core/discov"

type Config struct {
	EdgeId   int64
	ListenOn string
	Etcd     discov.EtcdConf
	Kq       struct {
		Brokers []string
		Topic   string
		GroupID string
	}
}
