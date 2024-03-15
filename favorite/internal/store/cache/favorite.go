package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// ensure this behavior is atomic
var releaseLockLua = redis.NewScript(`
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`)

func (m *Manager) ReleaseLock(k, v string) int {
	keys := []string{k}
	values := []string{v}
	val, err := releaseLockLua.Run(context.Background(), m.handler, keys, values).Int()
	if err != nil {
		return -1
	}
	return val
}

func (m *Manager) GetLock(k, v string, d time.Duration) bool {
	val, err := m.handler.SetNX(context.Background(), k, v, d).Result()
	if err != nil {
		return false
	}
	return val
}

func (m *Manager) GetExpiredKeys(keys []string, ttlChecker func(x, ttl time.Duration) bool, TTL time.Duration) []string {
	if len(keys) == 0 {
		return nil
	}
	ctx := context.Background()
	cmder, err := m.handler.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < len(keys); i++ {
			pipe.TTL(ctx, keys[i])
		}
		return nil
	})
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	ans := make([]string, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		if x := cmder[i].(*redis.DurationCmd).Val(); ttlChecker(x, TTL) {
			ans = append(ans, keys[i])
		}
	}
	return ans
}

func (m *Manager) CheckExpire(k string) time.Duration {
	val, err := m.handler.TTL(context.Background(), k).Result()
	if err != nil {
		return -1
	}
	return val
}

func (m *Manager) GetKeys(p string) []string {
	val, err := m.handler.Keys(context.Background(), p).Result()
	if err != nil {
		return nil
	}
	return val
}

// Pipeline in go-redis v8 is thread safe implemented with lock, but removed in v9
func (m *Manager) AddSetWithExpire(expire time.Duration, k string, v ...string) int64 {
	if len(v) == 0 {
		return -1
	}
	ctx := context.Background()
	var cmd *redis.IntCmd
	if _, err := m.handler.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		cmd = pipe.SAdd(ctx, k, v)
		pipe.Expire(ctx, k, expire)
		return nil
	}); err != nil {
		log.Println(err.Error())
		return -1
	}
	return cmd.Val()
}

func (m *Manager) RemoveSetValue(k, v string) int64 {
	ctx := context.Background()
	val, err := m.handler.SRem(ctx, k, v).Result()
	if err != nil {
		return -1
	}
	return val
}

func (m *Manager) DelKey(k string) int64 {
	ctx := context.Background()
	val, err := m.handler.Del(ctx, k).Result()
	if err != nil {
		return -1
	}
	return val
}

func (m *Manager) GetSetMember(k string) []string {
	val, err := m.handler.SMembers(context.Background(), k).Result()
	if err != nil {
		return nil
	}
	return val
}
