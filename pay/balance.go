package main

import (
	log "code.google.com/p/log4go"
	"errors"
	"fmt"
	"github.com/Terry-Mao/goim/libs/hash/ketama"
	inet "github.com/Terry-Mao/goim/libs/net"
	rpc "github.com/Terry-Mao/protorpc"
	proto "github.com/qinlodestar/pay-balance/proto/pay"
)

var (
	balanceServiceMap = map[string]**rpc.Client{}
	balanceRing       *ketama.HashRing
)

const (
	BalanceServicePush = "BalanceRpc.Push"
)

func initBalance() (err error) {
	var (
		network, addr string
	)
	balanceRing = ketama.NewRing(ketama.Base)
	for serverId, addrs := range Conf.BalanceRPCAddrs {
		// WARN r must every recycle changed for reconnect
		var (
			r          *rpc.Client
			routerQuit = make(chan struct{}, 1)
		)
		if network, addr, err = inet.ParseNetwork(addrs); err != nil {
			log.Error("inet.ParseNetwork() error(%v)", err)
			return
		}
		r, err = rpc.Dial(network, addr)
		if err != nil {
			log.Error("rpc.Dial(\"%s\", \"%s\") error(%v)", network, addr, err)
		}
		go rpc.Reconnect(&r, routerQuit, network, addr)
		log.Debug("router rpc addr:%s connect", addr)
		balanceServiceMap[serverId] = &r
		balanceRing.AddNode(serverId, 1)
	}
	balanceRing.Bake()
	return
}

func pushMsg(userId int64, orderId string, money float64) (err error) {
	client, err := getServerByUserId(userId)
	if err != nil {
		return
	}
	arg := &proto.MsgArg{UserId: userId, OrderId: orderId, Money: money}
	reply := &proto.NoReply{}
	if err = client.Call(BalanceServicePush, arg, reply); err != nil {
		return err
	}
	return
}

func getServerByUserId(userId int64) (*rpc.Client, error) {
	serverId := balanceRing.Hash(fmt.Sprintf("%d", userId))
	log.Debug(serverId)
	err := errors.New("get server wrong")
	if client, ok := balanceServiceMap[serverId]; !ok || *client == nil {
		return nil, err
	} else {
		return *client, nil
	}
}
