package redis

import (
	Config "github.com/tsukiamaoto/book-crawler-go/config"

	"context"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	RDB *redis.Client
}

var ctx = context.Background()

func New() *Redis {
	config := Config.LoadConfig()
	RDB := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	return &Redis{
		RDB: RDB,
	}
}

func (r *Redis) ConnectRDB() {
	_, err := r.RDB.Ping(ctx).Result()
	if err != nil {
		log.Error(err)
	}
}

func (r *Redis) Get(key string) string {
	value, err := r.RDB.Get(ctx, key).Result()
	if err != nil {
		log.Error(err)
	}
	return value
}

func (r *Redis) Set(key string, value interface{}) {
	_, err := r.RDB.Set(ctx, key, value, 0).Result()
	if err != nil {
		log.Error(err)
	}
}

func (r *Redis) GetEx(key string, expiration time.Duration) string {
	value, err := r.RDB.GetEx(ctx, key, expiration).Result()
	if err != nil {
		log.Error(err)
	}
	return value
}

func (r *Redis) SetEx(key string, value interface{}, expiration time.Duration) {
	_, err := r.RDB.Set(ctx, key, value, expiration).Result()
	if err != nil {
		log.Error(err)
	}
}

func (r *Redis) Exists(key string) bool {
	ok, err := r.RDB.Exists(ctx, key).Result()
	if err != nil {
		log.Error(err)
	}

	return ok == 1
}
