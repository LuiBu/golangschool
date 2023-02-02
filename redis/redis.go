package redis

import (
	"flag"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

// 声明一些全局变量
var (
	Conn          redis.Conn
	Pool          *redis.Pool
	redisServer   = flag.String("redisServer", "127.0.0.1:6379", "")
	redisPassword = flag.String("redisPassword", "liburdi/go-blogs", "")
)

// 初始化一个pool
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				defer c.Close()
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
func init() {
	flag.Parse()
	Pool = newPool(*redisServer, *redisPassword)
}
func Set(key, value string) interface{} {
	Conn = Pool.Get()
	defer Conn.Close()
	//redis操作
	v, err := Conn.Do("SET", key, value)
	if err != nil {
		fmt.Println(err)
	}
	//正常情况情况下返回ok,已设置返回nil,不需要多返回值
	return v
}
func Get(key string) (string, error) {
	Conn = Pool.Get()
	defer Conn.Close()
	v, err := redis.String(Conn.Do("GET", key))
	if err != nil {
		fmt.Println(err)
	}
	return v, err
}
