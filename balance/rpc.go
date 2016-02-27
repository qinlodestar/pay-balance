package main

import (
	log "code.google.com/p/log4go"
	inet "github.com/Terry-Mao/goim/libs/net"
	rpc "github.com/Terry-Mao/protorpc"
	proto "github.com/qinlodestar/pay-balance/proto/pay"
	"net"
)

type BalanceRpc struct {
}

func initRpc() (err error) {
	var (
		network, addr string
		c             = &BalanceRpc{}
	)
	rpc.Register(c)
	if network, addr, err = inet.ParseNetwork(Conf.RPCPushAddrs); err != nil {
		log.Error("inet.ParseNetwork() error(%v)", err)
		return
	}
	go rpcListen(network, addr)
	return
}

func rpcListen(network, addr string) {
	log.Debug("network=%s\tadd=%s", network, addr)
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Error("net.Listen(\"%s\", \"%s\") error(%v)", network, addr, err)
		panic(err)
	}
	// if process exit, then close the rpc addr
	defer func() {
		log.Info("listen rpc: \"%s\" close", addr)
		if err := l.Close(); err != nil {
			log.Error("listener.Close() error(%v)", err)
		}
	}()
	rpc.Accept(l)
}

func (this *BalanceRpc) Push(arg *proto.MsgArg, reply *proto.NoReply) (err error) {
	println("1111")
	return
}
