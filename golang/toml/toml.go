package toml

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/BurntSushi/toml"
	"github.com/EAHITechnology/inf/golang/emysql"
	"github.com/EAHITechnology/inf/golang/enet"
	"github.com/EAHITechnology/inf/golang/eredis"
	"github.com/EAHITechnology/inf/golang/log"
	"github.com/EAHITechnology/inf/golang/net_config"
)

type TomlConfig struct {
	ServerConfig    ServerConfigInfo     `toml:"server"`
	LogConfig       LogConfigInfo        `toml:"log"`
	DatabaseConfigs []DatabaseConfigInfo `toml:"database,omitempty"`
	RedisConfigs    []RedisConfigInfo    `toml:"redis,omitempty"`
	RpcNetConfigs   []RpcNetConfigInfo   `toml:"rpc_server_client"`
	HttpNetConfigs  []HttpNetConfigInfo  `toml:"http_server_client"`
}

func NewTomlConfig() *TomlConfig {
	return &TomlConfig{}
}

/*
Read函数需要传入一个toml文件的路径和一个接收对象，接收对象可以帮助返回配置文件中的业务信息
*/
func (t *TomlConfig) Read(path string, dst interface{}) error {
	if _, err := toml.DecodeFile(path, t); err != nil {
		fmt.Println(err)
		return err
	}
	if _, err := toml.DecodeFile(path, dst); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (t *TomlConfig) InitConfig(ctx context.Context) error {
	//init log
	if err := t.initLogger(t.LogConfig.Logpath); err != nil {
		return err
	}
	//init discover
	// if err := nw.Init(); err != nil {
	// 	return err
	// }
	//init base config
	enet.SetRpcLisPort(t.ServerConfig.RpcPort)
	enet.SetHttpLisPort(t.ServerConfig.HttpPort)
	//init mysql client
	if t.DatabaseConfigs != nil || len(t.DatabaseConfigs) != 0 {
		if err := emysql.NewMysql(t.mixMysql()); err != nil {
			panic("InitConfig NewMysql %s" + err.Error())
		}
		fmt.Println("InitConfig NewMysql success")
	}
	//init redis client
	if t.RedisConfigs != nil || len(t.RedisConfigs) != 0 {
		if err := eredis.InitRedis(t.mixRedis()); err != nil {
			panic("InitConfig InitRedis %s" + err.Error())
		}
		log.Info("InitConfig InitRedis success")
	}
	//init net config
	if err := net_config.InitClientRpc(ctx, t.mixRpc()); err != nil {
		return err
	}
	if err := net_config.InitClientHttp(ctx, t.mixHttp()); err != nil {
		return err
	}
	return nil
}

func (t *TomlConfig) mixMysql() []emysql.MConfigInfo {

	mcs := []emysql.MConfigInfo{}
	for _, m := range t.DatabaseConfigs {
		if m.MasterIp == "" {
			panic("mixMysql splitConnect MasterIp nil")
		}
		mc := emysql.MConfigInfo{
			Name:            m.Name,
			Username:        m.Username,
			MasterIp:        m.MasterIp,
			SlaveIp:         m.SlaveIp,
			Password:        m.Password,
			Database:        m.Database,
			Charset:         m.Charset,
			ParseTime:       m.ParseTime,
			Loc:             m.Loc,
			ReadTimeout:     m.ReadTimeout,
			MaxIdleConns:    m.MaxIdleConns,
			MaxOpenConns:    m.MaxOpenConns,
			ConnMaxLifetime: m.ConnMaxLifetime,
			DiscoverFlag:    m.DiscoverFlag,
		}
		mcs = append(mcs, mc)
	}
	return mcs
}

func (t *TomlConfig) mixRedis() []eredis.RedisInfo {
	redisInfos := []eredis.RedisInfo{}
	for _, r := range t.RedisConfigs {
		redisInfo := eredis.RedisInfo{
			RedisName:      r.Name,
			Addr:           r.Addr,
			MaxIdle:        r.MaxIdle,
			MaxActive:      r.MaxActive,
			IdleTimeout:    r.IdleTimeout,
			ReadTimeout:    r.ReadTimeout,
			WriteTimeout:   r.WriteTimeout,
			ConnectTimeout: r.ConnectTimeout,
			Password:       r.Password,
			Wait:           r.Wait,
		}
		redisInfos = append(redisInfos, redisInfo)
	}
	return redisInfos
}

func (t *TomlConfig) mixRpc() []net_config.RpcConnConfig {
	rs := []net_config.RpcConnConfig{}
	for _, r := range t.RpcNetConfigs {
		rpcConnConfig := net_config.RpcConnConfig{
			ServiceName:   r.ServiceName,
			Proto:         r.Proto,
			EndpointsFrom: r.EndpointsFrom,
			Addr:          r.Addr,
			Balancetype:   r.Balancetype,
			DialTimeout:   r.DialTimeout,
			ReadTimeout:   r.ReadTimeout,
			RetryTimes:    r.RetryTimes,
			MaxSize:       r.MaxSize,
		}
		rs = append(rs, rpcConnConfig)
	}
	return rs
}

func (t *TomlConfig) mixHttp() []net_config.HttpConnConfig {
	ht := []net_config.HttpConnConfig{}
	for _, h := range t.HttpNetConfigs {
		httpConnConfig := net_config.HttpConnConfig{
			ServiceName:   h.ServiceName,
			Proto:         h.Proto,
			EndpointsFrom: h.EndpointsFrom,
			Addr:          h.Addr,
			Balancetype:   h.Balancetype,
			ReadTimeout:   h.ReadTimeout,
			RetryTimes:    h.RetryTimes,
		}
		ht = append(ht, httpConnConfig)
	}
	return ht
}

func (t *TomlConfig) initLogger(path string) error {
	if err := createDir(path); err != nil {
		return err
	}
	log.Init(t.ServerConfig.ServiceName, path)
	return nil
}
