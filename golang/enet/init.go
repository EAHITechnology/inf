package enet

func SetHttpNetManager(serviceName string, ht *Ht) {
	httpMap[serviceName] = ht
}

func SetRpcNetManager(serviceName string, rp *Rp) {
	rpcMap[serviceName] = rp
}
