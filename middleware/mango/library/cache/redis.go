package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/x554462/gin-example/middleware/mango/library/conf"
	"sync"
	"time"
)

var redisDbOnce sync.Map
var redisClientMap sync.Map

type RedisClient struct {
	client *redis.Client
}

func NewRedis(db int) *RedisClient {
	var once sync.Once
	value, _ := redisDbOnce.LoadOrStore(db, &once)
	value.(*sync.Once).Do(func() {
		c := conf.RedisConf
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
			Password: c.Password,
			DB:       db,
		})
		redisClientMap.Store(db, &RedisClient{client: redisClient})
	})
	client, _ := redisClientMap.Load(db)
	return client.(*RedisClient)
}

func (r *RedisClient) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(key, expiration).Err()
}

func (r *RedisClient) HGet(key, field string) (v interface{}) {
	r.client.HGet(key, field).Scan(&v)
	return
}

func (r *RedisClient) HGetString(key, field string) string {
	return r.client.HGet(key, field).String()
}

func (r *RedisClient) HSet(key, field string, value interface{}) error {
	return r.client.HSet(key, field, value).Err()
}

func (r *RedisClient) HExists(key, field string) bool {
	if exists, err := r.client.HExists(key, field).Result(); err == nil && exists {
		return true
	}
	return false
}

func (r *RedisClient) HDel(key string, field ...string) error {
	return r.client.HDel(key, field...).Err()
}

func (r *RedisClient) Del(keys ...string) error {
	r.client.Context()
	return r.client.Del(keys...).Err()
}
