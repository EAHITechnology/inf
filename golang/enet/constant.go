package enet

const (
	DESIGNATED = "designated"
	ROUNDROBIN = "roundrobin"
	RANDOM     = "random"
)

const (
	CONTENT_TYPE      = "Content-type"
	CONTENT_TYPE_JSON = "application/json"
)

const (
	LOGID = "log_id"
)

var rpcMap = make(map[string]*Rp)
var httpMap = make(map[string]*Ht)

var RPC_PORT string = ""
var HTTP_PORT string = ""
