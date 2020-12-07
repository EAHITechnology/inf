package toml

type ServerConfigInfo struct {
	ServiceName string `toml:"service_name"`
	HttpPort    string `toml:"http_port"`
	RpcPort     string `toml:"rpc_port"`
}

type LogConfigInfo struct {
	Logpath string `toml:"logpath"`
}

type DatabaseConfigInfo struct {
	Name            string `toml:"name"`
	Username        string `toml:"username"`
	MasterIp        string `toml:"master"`
	SlaveIp         string `toml:"slave"`
	Password        string `toml:"password"`
	Database        string `toml:"database"`
	Charset         string `toml:"charset"`
	ParseTime       string `toml:"parseTime"`
	Loc             string `toml:"loc"`
	ReadTimeout     string `toml:"readTimeout"`
	MaxIdleConns    int    `toml:"maxIdleConns"`
	MaxOpenConns    int    `toml:"maxOpenConns"`
	ConnMaxLifetime int    `toml:"connMaxLifetime"`
	DiscoverFlag    bool   `toml:"discover_flag"`
}

type mConfigInfo struct {
	Name         string
	Master       string
	Slave        string
	MaxIdleConns int
	MaxOpenConns int
}

type RedisConfigInfo struct {
	Name           string `toml:"name"`
	Addr           string `toml:"addr"`
	Password       string `toml:"password"`
	MaxIdle        int    `toml:"max_idle"`
	IdleTimeout    int64  `toml:"max_idletimeout"`
	MaxActive      int    `toml:"max_active"`
	ReadTimeout    int64  `toml:"read_timeout"`
	WriteTimeout   int64  `toml:"write_timeout"`
	SlowTime       int64  `toml:"slow_time"`
	ConnectTimeout int64  `toml:"connect_time"`
	Wait           int8   `toml:"wait"`
}

type RpcNetConfigInfo struct {
	ServiceName   string   `toml:"service_name"`
	Proto         string   `toml:"proto"`
	EndpointsFrom string   `toml:"endpoints_from"`
	Addr          []string `toml:"addr"`
	Balancetype   string   `toml:"balancetype"`
	DialTimeout   int      `toml:"dial_timeout"`
	ReadTimeout   int      `toml:"read_timeout"`
	RetryTimes    int      `toml:"retry_times"`
	MaxSize       int      `toml:"max_size"`
}

type HttpNetConfigInfo struct {
	ServiceName   string   `toml:"service_name"`
	Proto         string   `toml:"proto"`
	EndpointsFrom string   `toml:"endpoints_from"`
	Addr          []string `toml:"addr"`
	Balancetype   string   `toml:"balancetype"`
	ReadTimeout   int      `toml:"read_timeout"`
	RetryTimes    int      `toml:"retry_times"`
}
