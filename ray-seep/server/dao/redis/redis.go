package redis

import (
	"github.com/go-redis/redis"
)

type Client struct {
	redis.BitCount
}
