package redis

import (
	"context"
	"time"
)

const (
	// 失败重试次数
	defaultMaxRetries = 3
	// 连接池容量
	defaultRedisPoolSize = 15
	// 空闲连接数
	defaultRedisMinIdleConn = 10
	// 连接超时
	defaultDialTimeout = 5 * time.Second
	// 写超时
	defaultWriteTimeout = 3 * time.Second
	// 读超时
	defaultReadTimeout = 3 * time.Second
)

type RedisOption struct {
	Addr string
	Pass string
	DB   int
}

type cacheRepo struct {
	Opt    RedisOption
	Client *redis.Client
}

func NewRedis(o RedisOption) (*cacheRepo, error) {
	c := &cacheRepo{Opt: o}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *cacheRepo) connect() error {
	c.Client = redis.NewClient(&redis.Options{
		Addr:         c.Opt.Addr,
		Password:     c.Opt.Pass,
		DB:           c.Opt.DB,
		MaxRetries:   defaultMaxRetries,
		PoolSize:     defaultRedisPoolSize,
		MinIdleConns: defaultRedisMinIdleConn,
		DialTimeout:  defaultDialTimeout,
		WriteTimeout: defaultWriteTimeout,
		ReadTimeout:  defaultReadTimeout,
	})
	return c.Client.Ping(context.Background()).Err()
}
