package deps

import "github.com/redis/go-redis/v9"

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedis(cfg RedisConfig) *redis.Client {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return redisDB
}
