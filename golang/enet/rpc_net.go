package enet

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Rp struct {
	serviceName   string
	addr          []string
	proto         string
	endpointsFrom string
	balancetype   string
	dialTimeout   int
	readTimeout   int
	retryTimes    int
	maxSize       int
	addrIdx       int
}

func NewRp() *Rp {
	return &Rp{}
}

func (r *Rp) SetServiceName(name string) *Rp {
	r.serviceName = name
	return r
}

func (r *Rp) SetAddr(addr []string) *Rp {
	r.addr = addr
	return r
}

func (r *Rp) SetEndpointsFrom(endpointsFrom string) *Rp {
	r.endpointsFrom = endpointsFrom
	return r
}

func (r *Rp) SetBalancetype(balancetype string) *Rp {
	r.balancetype = balancetype
	return r
}

func (r *Rp) SetDialTimeout(dialTimeout int) *Rp {
	r.dialTimeout = dialTimeout
	return r
}

func (r *Rp) SetReadTimeout(readTimeout int) *Rp {
	r.readTimeout = readTimeout
	return r
}

func (r *Rp) SetRetryTimes(retryTimes int) *Rp {
	r.retryTimes = retryTimes
	return r
}

func (r *Rp) SetMaxSize(maxSize int) *Rp {
	r.maxSize = maxSize
	return r
}

func (r *Rp) getMaxSize() int {
	if r.maxSize < 1024 {
		r.maxSize = 1024 * 1024 * 4
	}
	return r.maxSize
}

/*添加请求地址*/
func (r *Rp) getAddr() (string, error) {
	var addr string = ""
	switch r.endpointsFrom {
	case "etcd":
	case "k8s":
	case "consoul":
	case "ip":
		lenth := len(r.addr)
		if lenth == 0 {
			return "", errors.New("addr is not init")
		}
		if r.balancetype == RANDOM {
			addr = r.addr[rand.Intn(lenth)]
		}
		if r.balancetype == ROUNDROBIN {
			if r.addrIdx >= lenth {
				r.addrIdx = 0
			}
			addr = r.addr[r.addrIdx]
			r.addrIdx = (r.addrIdx + 1) % lenth
		}
	default:
		return "", errors.New("error endpointsFrom")
	}
	return addr, nil
}

/*获得链接*/
func NewRpcConn(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	r, ok := rpcMap[serviceName]
	if !ok {
		return nil, errors.New("serviceName not init")
	}
	addr, err := r.getAddr()
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(
		addr, grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Duration(r.dialTimeout)*time.Millisecond),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(r.getMaxSize())))
	if err != nil {
		return conn, errors.New(fmt.Sprintf("can not connect to %s error %v", addr, err))
	}
	return conn, nil
}

/*将链接放回链接池*/
func CloseRpcConn(conn *grpc.ClientConn) {
	if conn != nil {
		conn.Close()
	}
}

/*设置监听端口*/
func SetRpcLisPort(rpcPort string) {
	RPC_PORT = rpcPort
}

/*初始化grpc监听*/
func IniRpctLis() (*grpc.Server, *net.Listener, error) {
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", ":"+RPC_PORT)
	if err != nil {
		fmt.Println("Listen IniRpctLis error:" + err.Error())
		panic("Listen IniRpctLis err:" + err.Error())
	}
	fmt.Println("IniRpctLis port: " + RPC_PORT + " success")
	return s, &lis, nil
}

/*注册反射服务*/
func ReflectionRegister(s *grpc.Server) {
	reflection.Register(s)
}
