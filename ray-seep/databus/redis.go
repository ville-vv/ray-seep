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
	return c.rds.HSetNX("login_token_", user, token).Err()
}

// 更新到期时间
func (c *RedisClient) UpdateTokenTTl(user string, id int64) error {
	return c.rds.Expire(fmt.Sprintf("%s_%d_token", user, id), time.Second*60).Err()
}

func (c *RedisClient) DelUserToken(connID int64, user string, isDelKeys bool) error {
	return c.rds.HDel("login_token_", user).Err()
}

func (c *RedisClient) GetUserToken(connID int64, user string) string {
	return c.rds.HGet("login_token_", user).Val()
}
