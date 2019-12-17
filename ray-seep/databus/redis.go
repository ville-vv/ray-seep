package databus

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"vilgo/vredis"
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
func (c *RedisClient) SetUserToken(user string, id int64, token string) error {
	fmt.Println(strconv.FormatInt(id, 10))
	return c.rds.HSetNX(user, strconv.FormatInt(id, 10), token).Err()
}

// 更新到期时间
func (c *RedisClient) UpdateTokenTTl(user string, id int64) error {
	return c.rds.Expire(fmt.Sprintf("%s_%d_token", user, id), time.Second*60).Err()
}

func (c *RedisClient) DelUserToken(user string, id int64) error {
	return c.rds.HDel(user, strconv.FormatInt(id, 10)).Err()
}
