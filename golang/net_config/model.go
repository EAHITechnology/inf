package net_config

type RpcConnConfig struct {
	ServiceName   string
	Proto         string
	Addr          []string
	EndpointsFrom string
	Balancetype   string
	DialTimeout   int
	ReadTimeout   int
	RetryTimes    int
	MaxSize       int
}

type HttpConnConfig struct {
	ServiceName   string
	Proto         string
	EndpointsFrom string
	Addr          []string
	Balancetype   string
	ReadTimeout   int
	RetryTimes    int
}
