package nosql

import (
	"context"
	"encoding/json"
	"time"

	"github.com/KAMIENDER/golang-scaffold/infra/config"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// Redis implement r *Redis NoSQLDB
type Redis struct {
	client *redis.Client
}

var redisClient NoSQLDB = new(Redis)

func NewRedis(conf *config.Config) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisConf.Addr,
		Password: conf.RedisConf.Password, // 没有密码，默认值
		DB:       conf.RedisConf.DB,       // 默认DB 0
	})
	return &Redis{
		client: rdb,
	}
}

func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	cmd := r.client.Set(ctx, key, value, expiration)
	return errors.Wrap(cmd.Err(), "")
}

func (r *Redis) Get(ctx context.Context, key string, obj any) (bool, error) {
	cmd := r.client.Get(ctx, key)
	str, err := cmd.Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, errors.Wrap(json.Unmarshal([]byte(str), obj), "")
}

func (r *Redis) Del(ctx context.Context, key string) error {
	cmd := r.client.Del(ctx, key)
	return errors.Wrap(cmd.Err(), "")
}
