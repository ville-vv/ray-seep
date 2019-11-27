package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Client struct {
	rds *redis.Client
}

func NewClient(cli *redis.Client) *Client {
	return &Client{
		rds: cli,
	}
}

func (c *Client) SetUserToken(user string, id int, token string) error {
	key := fmt.Sprintf("%s_%d_token", user, id)
	cmd := c.rds.Set(key, token, 0)
	return cmd.Err()
}
