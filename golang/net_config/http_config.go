package net_config

import (
	"fmt"

	"github.com/EAHITechnology/inf/golang/enet"
	"golang.org/x/net/context"
)

func InitClientHttp(ctx context.Context, httpConnConfig []HttpConnConfig) error {
	if len(httpConnConfig) == 0 {
		return nil
	}
	for _, h := range httpConnConfig {
		if err := checkHttpParam(&h); err != nil {
			Http = nil
			return err
		}
		Http[h.ServiceName] = h
		ht := enet.NewHt()
		ht.SetAddr(h.Addr).SetBalancetype(h.Balancetype).
			SetEndpointsFrom(h.EndpointsFrom).SetProto(h.Proto).
			SetReadTimeout(h.ReadTimeout).SetRetryTimes(h.RetryTimes).
			SetServiceName(h.ServiceName)
		enet.SetHttpNetManager(h.ServiceName, ht)
	}
	fmt.Printf("InitClientHttp http:%v", Http)
	return nil
}

func checkHttpParam(h *HttpConnConfig) error {
	if h.ServiceName == "" {
		return fmt.Errorf("service name nil")
	}
	if h.Proto == "" {
		return fmt.Errorf("service proto nil")
	}
	if h.EndpointsFrom == "" {
		h.EndpointsFrom = "sfns"
		return nil
	}
	if h.Balancetype == "" {
		h.Balancetype = "random"
	}
	if h.ReadTimeout == 0 {
		h.ReadTimeout = 1000
	}
	return nil
}
