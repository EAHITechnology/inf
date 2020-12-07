package eredis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	pool *redis.Pool
}

func checkConfig(r *RedisInfo) error {
	if r.RedisName == "" {
		return fmt.Errorf("no RedisName")
	}
	if r.Addr == "" {
		return fmt.Errorf("no Addr")
	}
	if r.MaxIdle == 0 {
		r.MaxIdle = 50
	}
	if r.MaxActive == 0 {
		r.MaxActive = 100
	}
	if r.IdleTimeout == 0 {
		r.ReadTimeout = 300
	}
	if r.ReadTimeout == 0 {
		r.ReadTimeout = 500
	}
	if r.WriteTimeout == 0 {
		r.WriteTimeout = 500
	}
	if r.ConnectTimeout == 0 {
		r.ConnectTimeout = 500
	}
	return nil
}

func InitRedis(redisInfos []RedisInfo) error {
	for _, i := range redisInfos {
		red := new(Redis)
		if err := checkConfig(&i); err != nil {
			continue
		}
		wait := false
		if i.Wait == 1 {
			wait = true
		}
		red.pool = &redis.Pool{
			MaxIdle:     i.MaxIdle,
			MaxActive:   i.MaxActive,
			IdleTimeout: time.Duration(i.IdleTimeout) * time.Second,
			Wait:        wait,
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					i.Addr,
					redis.DialPassword(i.Password),
					redis.DialReadTimeout(time.Duration(i.ReadTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(i.WriteTimeout)*time.Millisecond),
					redis.DialConnectTimeout(time.Duration(i.ConnectTimeout)*time.Millisecond),
					redis.DialDatabase(0),
				)
			},
		}
		redisclient[i.RedisName] = red
	}
	return nil
}

func (this *Redis) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := this.pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}

func (this *Redis) Del(args ...interface{}) (count int, err error) {
	con := this.pool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()
	var reply interface{}
	reply, err = con.Do("Del", args...)
	if err != nil {
		return
	}
	count = reply.(int)
	return
}

func getClients() map[string]*Redis {
	return redisclient
}

func GetClient(redisName string) (*Redis, error) {
	rediss := getClients()
	r, ok := rediss[redisName]
	if !ok {
		return nil, fmt.Errorf("redis not init")
	}
	return r, nil
}
