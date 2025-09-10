package config

import (
	"auth-service/db"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/xinyi-chong/common-lib/logger"
	redisclient "github.com/xinyi-chong/common-lib/redis"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port" validate:"required"`
	} `mapstructure:"server"`

	Postgres struct {
		db.Config `mapstructure:",squash"`
	} `mapstructure:"auth_postgres"`

	Redis struct {
		redisclient.Config `mapstructure:",squash"`
	} `mapstructure:"redis"`

	JWT struct {
		SecretKey       string        `mapstructure:"secret_key" validate:"required"`
		AccessDuration  time.Duration `mapstructure:"access_duration"`
		RefreshDuration time.Duration `mapstructure:"refresh_duration"`
	} `mapstructure:"jwt"`
}

func (c *Config) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			fmt.Println("Validation errors:")
			for _, e := range validationErrors {
				fmt.Printf("Field: %s, Tag: %s, Value: %v\n", e.Field(), e.Tag(), e.Value())
			}
		}
		return err
	}
	return nil
}

func Load() (*Config, error) {
	v := setupViper()

	if err := v.ReadInConfig(); err != nil {
		logger.Error("Load: ReadInConfig: ", zap.Error(err))
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		logger.Error("Load: Unmarshal config: ", zap.Error(err))
		return nil, err
	}
	logger.Debug("Load: config: ", zap.Any("config", config))

	if err := config.Validate(); err != nil {
		logger.Error("Load: Validate: ", zap.Error(err))
		return nil, err
	}

	return &config, nil
}

func setupViper() *viper.Viper {
	v := viper.New()

	v.SetConfigName("auth-service")
	v.SetConfigType("yaml")
	v.AddConfigPath("/app/configs")
	v.AddConfigPath("./configs") // For local development

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v
}
