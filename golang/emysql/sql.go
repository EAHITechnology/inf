package emysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type mClient struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

/*
NewMysql方法提供配置文件读取mysql配置，不允许业务侧直接添加，所有数据库创建应该在初始化的时候都写在配置中
*/
func NewMysql(mcs []MConfigInfo) error {
	for _, m := range mcs {
		if m.Name == "" {
			return ErrNameInvalid
		}
		if m.Username == "" {
			return ErrUserNameInvalid
		}
		if m.Password == "" {
			return ErrPassWordInvalid
		}
		if m.MasterIp == "" {
			return ErrMasterInvalid
		}
		if m.MaxIdleConns == 0 {
			m.MaxIdleConns = MAX_IDLE_CONNS
		}
		if m.MaxOpenConns == 0 {
			m.MaxOpenConns = MAX_OPEN_CONNS
		}
		if m.ConnMaxLifetime == 0 {
			m.ConnMaxLifetime = CONN_MAX_LiFE_TIME
		}
		if m.DiscoverFlag {
			//通过服务发现解析mysql地址
		}
		master := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s&readTimeout=%s",
			m.Username, m.Password,
			m.MasterIp, m.Database,
			m.Charset, m.ParseTime,
			m.Loc, m.ReadTimeout,
		)
		mc := &mClient{}
		masterdb, err := gorm.Open("mysql", master)
		if err != nil {
			return err
		}
		masterdb.DB().SetMaxIdleConns(m.MaxIdleConns)
		masterdb.DB().SetMaxOpenConns(m.MaxOpenConns)
		masterdb.DB().SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime) * time.Second)
		mc.Master = masterdb

		if m.SlaveIp != "" {
			slave := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s&readTimeout=%s",
				m.Username, m.Password,
				m.SlaveIp, m.Database,
				m.Charset, m.ParseTime,
				m.Loc, m.ReadTimeout,
			)
			slavedb, err := gorm.Open("mysql", slave)
			if err != nil {
				return err
			}
			slavedb.DB().SetMaxIdleConns(m.MaxIdleConns)
			slavedb.DB().SetMaxOpenConns(m.MaxOpenConns)
			slavedb.DB().SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime) * time.Second)
			mc.Slave = slavedb
		}
		dbclient[m.Name] = mc
	}
	return nil
}

func CloseMysql() error {
	for _, v := range dbclient {
		if v.Master != nil {
			v.Master.Close()
		}
		if v.Slave != nil {
			v.Slave.Close()
		}
	}
	return nil
}

func getClients() map[string]*mClient {
	return dbclient
}

func GetClient(dbname string) (*mClient, error) {
	dbs := getClients()
	db, ok := dbs[dbname]
	if !ok {
		return nil, fmt.Errorf("db not init")
	}
	return db, nil
}
