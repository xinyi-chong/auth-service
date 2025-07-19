package redisclient

import (
	"auth-service/pkg/logger"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	instance *redis.Client
	once     sync.Once
)

type Config struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DB       int    `mapstructure:"db" validate:"gte=0,lte=15"`
}

func Client() (*redis.Client, error) {
	if instance == nil {
		return nil, errors.New("redis client not initialized")
	}
	return instance, nil
}

func Init(cfg Config) (*redis.Client, error) {
	var err error
	once.Do(func() {
		instance = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		logger.Debug("Init: Connecting Redis client...",
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
			zap.Int("db", cfg.DB))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = instance.Ping(ctx).Err()
		if err != nil {
			logger.Error("Init: Failed to connect to Redis", zap.Error(err))
			instance = nil
			return
		}

		logger.Info("Init: Redis client connected")
	})

	return instance, err
}

func Close() error {
	if instance != nil {
		return instance.Close()
	}
	return nil
}

func Get(c context.Context, key string) (string, error) {
	client, err := Client()
	if err != nil {
		return "", err
	}
	return client.Get(c, key).Result()
}

func Set(c context.Context, key string, value interface{}, expiration time.Duration) error {
	client, err := Client()
	if err != nil {
		return err
	}
	err = client.Set(c, key, value, expiration).Err()
	return err
}

func Exists(c context.Context, key string) (bool, error) {
	client, err := Client()
	if err != nil {
		return false, err
	}
	result, err := client.Exists(c, key).Result()
	return result == 1, err
}
