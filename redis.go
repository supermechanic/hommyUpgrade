package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
)

//RedisClient redis连接池对象
var RedisClient *redis.Pool

func init() {
	host := Config.Redis.Base.Address + ":" + Config.Redis.Base.Port
	maxIdle := Config.Redis.MaxIdle
	maxActive := Config.Redis.MaxActive
	maxIdleTimeout := time.Duration(Config.Redis.IdleTimeout)
	timeout := time.Duration(Config.Redis.Timeout)

	// 建立连接池
	RedisClient = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: maxIdleTimeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", host,
				redis.DialConnectTimeout(timeout*time.Second),
				redis.DialReadTimeout(timeout*time.Second),
				redis.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			return con, nil
		},
	}
	fmt.Printf("%+v\n",RedisClient)
}
