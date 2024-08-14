package main

import (
	"flag"
	"log"
	"zeroim/core/protocol"
	"zeroim/core/socket"
	"zeroim/edge/internal/config"
	"zeroim/edge/internal/server/tcpserver"
	"zeroim/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("config", "../etc/edge-tcp.yaml", "the file of config")

func main() {
	flag.Parse()
	cfg := &config.Config{}
	err := conf.Load(*configFile, cfg)
	if err != nil {
		log.Printf("Config Load error: %v \n", err)
		return
	}
	svcCtx, err := svc.NewServiceContext(cfg)
	if err != nil {
		log.Printf("New ServiceContext error: %v \n", err)
		return
	}
	s := tcpserver.NewServer(svcCtx)
	s.RegistryHandle()

	server := socket.NewServer(svcCtx.Manager, uint64(cfg.EdgeId), protocol.ImProtocol{})
	defer server.Close()
	log.Printf("Listen on %s ...\n", cfg.ListenOn)

	server.ListenAndServe("tcp", cfg.ListenOn, nil)
}
