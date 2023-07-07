package local

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	pool        *redis.Pool
	ConnectTime time.Duration
	ReadTime    time.Duration
	WriteTime   time.Duration
	Password    string
}

func NewDefaultRedis() *Redis {
	return &Redis{
		ConnectTime: time.Second,
		ReadTime:    time.Second,
		WriteTime:   time.Second,
		Password:    "",
	}
}

func NewRedis(addr string) *Redis {
	return &Redis{}
}

func (r *Redis) Init() {
	r.pool = &redis.Pool{
		MaxIdle:     30,
		MaxActive:   30,
		IdleTimeout: time.Duration(time.Second),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "192.168.1.44:6377",
				redis.DialConnectTimeout(r.ConnectTime),
				redis.DialReadTimeout(r.ReadTime),
				redis.DialWriteTimeout(r.WriteTime),
				redis.DialPassword(r.Password),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func (r *Redis) Close() {
	r.pool.Close()
}
