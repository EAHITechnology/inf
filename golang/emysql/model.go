package emysql

type MConfigInfo struct {
	Name            string
	Username        string
	MasterIp        string
	SlaveIp         string
	Password        string
	Database        string
	Charset         string
	ParseTime       string
	Loc             string
	ReadTimeout     string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	DiscoverFlag    bool
}
