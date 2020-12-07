package net_config

var Rpc = make(map[string]RpcConnConfig)
var Http = make(map[string]HttpConnConfig)

const (
	defaultMaxReceiveMessageSize = 1024 * 1024 * 4
)
