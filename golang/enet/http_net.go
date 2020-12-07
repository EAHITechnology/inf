package enet

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

type Ht struct {
	serviceName   string
	addr          []string
	proto         string
	endpointsFrom string
	balancetype   string
	readTimeout   int
	retryTimes    int
	addrIdx       int
}

func NewHt() *Ht {
	return &Ht{}
}

/*设置服务名*/
func (h *Ht) SetServiceName(name string) *Ht {
	h.serviceName = name
	return h
}

/*设置服务地址,可用作兜底，可用作请求*/
func (h *Ht) SetAddr(addr []string) *Ht {
	h.addr = addr
	return h
}

/*设置服务发现方式*/
func (h *Ht) SetEndpointsFrom(endpointsFrom string) *Ht {
	h.endpointsFrom = endpointsFrom
	return h
}

/*设置负载均衡方式*/
func (h *Ht) SetBalancetype(balancetype string) *Ht {
	h.balancetype = balancetype
	return h
}

/*设置读超时*/
func (h *Ht) SetReadTimeout(readTimeout int) *Ht {
	h.readTimeout = readTimeout
	return h
}

/*设置重试次数*/
func (h *Ht) SetRetryTimes(retryTimes int) *Ht {
	h.retryTimes = retryTimes
	return h
}

/*设置协议类型*/
func (h *Ht) SetProto(proto string) *Ht {
	h.proto = proto
	return h
}

/*添加get参数*/
func (h *Ht) addParam(req *http.Request, params map[string]string) {
	if params != nil {
		q := req.URL.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
}

/*添加请求头部参数*/
func (h *Ht) addHeader(req *http.Request, headers map[string]string) {
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
}

/*添加请求地址*/
func (h *Ht) getAddr() (string, error) {
	var addr string = ""
	var err error
	switch h.endpointsFrom {
	case "etcd":
	case "k8s":
	case "consoul":
	case "ip":
		lenth := len(h.addr)
		if lenth == 0 {
			return "", errors.New("addr is not init")
		}
		if h.balancetype == RANDOM {
			addr = h.addr[rand.Intn(lenth)]
		}
		if h.balancetype == ROUNDROBIN {
			if h.addrIdx >= lenth {
				h.addrIdx = 0
			}
			addr = h.addr[h.addrIdx]
			h.addrIdx = (h.addrIdx + 1) % lenth
		}
	default:
		return "", errors.New("error endpointsFrom")
	}
	fullAddr, err := Write(h.proto, "://", addr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("getAddr Write err: %v", err))
	}
	return fullAddr, nil
}

/*基本get方法*/
func Get(ctx context.Context, service string, params map[string]string, headers map[string]string, url string) (*[]byte, int, error) {
	h, ok := httpMap[service]
	if !ok {
		return nil, 0, errors.New("service not init")
	}
	times := h.retryTimes
	for {
		times--
		addr, err := h.getAddr()
		if err != nil {
			if times < 0 {
				return nil, 0, err
			}
			continue
		}
		//构建请求
		req, err := http.NewRequest(http.MethodGet, addr+url, nil)
		if err != nil {
			if times < 0 {
				return nil, 0, errors.New(fmt.Sprintf("new request is fail err:%v", err))
			}
			continue
		}
		//添加url参数
		h.addParam(req, params)
		if logId := ctx.Value(LOGID); logId != nil {
			req.Header.Set(LOGID, logId.(string))
		}
		//添加头部
		h.addHeader(req, headers)
		//构建client配置
		client := &http.Client{Timeout: time.Duration(h.readTimeout) * time.Millisecond} //todo
		//发送请求
		resp, err := client.Do(req)
		if err != nil {
			if times < 0 {
				return nil, 0, err
			}
			continue
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, errors.New(fmt.Sprintf("http get ReadAll body failed: %v \n", err))
		}
		if resp.StatusCode != http.StatusOK {
			return &data, resp.StatusCode, nil
		}
		return &data, resp.StatusCode, nil
	}
}

/*基本post方法，从启动时配置文件加载配置*/
func Post(ctx context.Context, service string, body io.Reader, params map[string]string, headers map[string]string, url, contentType string) (*[]byte, int, error) {
	h, ok := httpMap[service]
	if !ok {
		return nil, 0, errors.New("service not init")
	}
	times := h.retryTimes
	for {
		times--
		//获得地址
		addr, err := h.getAddr()
		if err != nil {
			if times < 0 {
				return nil, 0, err
			}
			continue
		}
		//构造请求
		req, err := http.NewRequest(http.MethodPost, addr+url, body)
		if err != nil {
			if times < 0 {
				return nil, 0, errors.New(fmt.Sprintf("new request is fail err:%v", err))
			}
			continue
		}
		//添加头部
		if contentType == "" {
			contentType = CONTENT_TYPE_JSON
		}
		req.Header.Set(CONTENT_TYPE, contentType)
		h.addHeader(req, headers)
		//添加url
		h.addParam(req, params)
		if logId := ctx.Value(LOGID); logId != nil {
			req.Header.Set(LOGID, logId.(string))
		}
		//添加参数
		client := &http.Client{Timeout: time.Duration(h.readTimeout) * time.Millisecond} //需要区别各个阶段时间
		//发送请求
		resp, err := client.Do(req)
		if err != nil {
			if times < 0 {
				return nil, 0, err
			}
			continue
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, err
		}
		if resp.StatusCode != http.StatusOK {
			return &data, resp.StatusCode, nil
		}
		return &data, resp.StatusCode, nil
	}
}
