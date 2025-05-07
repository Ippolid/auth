package redis

import (
	"github.com/Ippolid/auth/internal/client/cache/redis"
	"github.com/Ippolid/auth/internal/repository"
)

type cache struct {
	cl redis.Client
}

// NewRedisCache конструктор для redis.
func NewRedisCache(cl redis.Client) repository.CacheInterface {
	return &cache{cl: cl}
}
