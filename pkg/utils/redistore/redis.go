package redistore

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
)

type RedisClusterPool struct {
	rc *redisc.Cluster
}

func NewRedisClusterPool(rc *redisc.Cluster) *RedisClusterPool {
	return &RedisClusterPool{rc: rc}
}

func (r *RedisClusterPool) Get() redis.Conn {
	return r.rc.Get()
}

func (r *RedisClusterPool) Close() error {
	return r.Close()
}
