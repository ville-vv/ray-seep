package databus

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/vilsongwei/vilgo/vredis"
	"time"
)

type RedisClient struct {
	rds *redis.Client
}

func NewRedisClient(cfg *vredis.RedisCnf) *RedisClient {
	rds := vredis.NewGoRedisDrive(cfg)
	if err := rds.Conn(); err != nil {
		panic(err)
	}
	if err := rds.GetRedis().Ping().Err(); err != nil {
		panic(err)
	}
	return &RedisClient{
		rds: rds.GetRedis(),
	}
}

func (c *RedisClient) Close() error {
	return c.rds.Close()
}

// 设置 Token
func (c *RedisClient) SetUserToken(connID int64, user string, token string) error {
	return c.rds.SetNX("rayseep::login::token::"+user, token, time.Second*60).Err()
}

// 更新到期时间
func (c *RedisClient) UpdateTokenTTl(user string, id int64) error {
	return c.rds.Expire(fmt.Sprintf("rayseep::login::token::"+user), time.Second*60).Err()
}

func (c *RedisClient) DelUserToken(connID int64, user string) error {
	return c.rds.Del("rayseep::login::token::" + user).Err()
}

func (c *RedisClient) GetUserToken(connID int64, user string) string {
	return c.rds.Get("rayseep::login::token::" + user).Val()
}
