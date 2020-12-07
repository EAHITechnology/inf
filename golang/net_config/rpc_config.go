package net_config

import (
	"fmt"

	"github.com/EAHITechnology/inf/golang/enet"
	"golang.org/x/net/context"
)

func InitClientRpc(ctx context.Context, rpcConnConfig []RpcConnConfig) error {
	if len(rpcConnConfig) == 0 {
		return nil
	}
	for _, r := range rpcConnConfig {
		if err := checkRpcParam(&r); err != nil {
			Rpc = nil
			return err
		}
		rp := enet.NewRp()
		rp.SetAddr(r.Addr).SetBalancetype(r.Balancetype).
			SetDialTimeout(r.DialTimeout).SetEndpointsFrom(r.EndpointsFrom).
			SetMaxSize(r.MaxSize).SetReadTimeout(r.ReadTimeout).
			SetRetryTimes(r.RetryTimes).SetServiceName(r.ServiceName)
		enet.SetRpcNetManager(r.ServiceName, rp)
	}
	fmt.Printf("InitClientRpc rpc:%+v", Rpc)
	return nil
}

func checkRpcParam(r *RpcConnConfig) error {
	if r.ServiceName == "" {
		return fmt.Errorf("service name nil")
	}
	if r.Proto == "" {
		return fmt.Errorf("service proto nil")
	}
	if r.EndpointsFrom == "" {
		return fmt.Errorf("endpoints from nil")
	}
	if r.Balancetype == "" {
		r.Balancetype = "roundrobin"
	}
	if r.DialTimeout == 0 {
		r.DialTimeout = 500
	}
	if r.ReadTimeout == 0 {
		r.ReadTimeout = 1000
	}
	if r.RetryTimes == 0 {
		r.RetryTimes = 1
	}
	if r.MaxSize == 0 {
		r.MaxSize = defaultMaxReceiveMessageSize
	}
	return nil
}
