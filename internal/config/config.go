package config

import (
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq" // postgres driver don`t delete
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Log  struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"log"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type DatabaseConfig struct {
	SQL struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"sql"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
}

type JWTConfig struct {
	User struct {
		AccessToken struct {
			SecretKey     string        `mapstructure:"secret_key"`
			TokenLifetime time.Duration `mapstructure:"token_lifetime"`
		} `mapstructure:"access_token"`
		RefreshToken struct {
			SecretKey     string        `mapstructure:"secret_key"`
			EncryptionKey string        `mapstructure:"encryption_key"`
			TokenLifetime time.Duration `mapstructure:"token_lifetime"`
		} `mapstructure:"refresh_token"`
	} `mapstructure:"user"`
	Service struct {
		SecretKey string `mapstructure:"secret_key"`
	} `mapstructure:"service"`
	Invites struct {
		SecretKey string `mapstructure:"secret_key"`
	} `mapstructure:"gov-invites"`
}

type SwaggerConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
	Port    string `mapstructure:"port"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Database DatabaseConfig `mapstructure:"database"`
	Swagger  SwaggerConfig  `mapstructure:"swagger"`
}

func LoadConfig() Config {
	configPath := os.Getenv("KV_VIPER_FILE")
	if configPath == "" {
		panic(fmt.Errorf("KV_VIPER_FILE env var is not set"))
	}
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %s", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("error unmarshalling config: %s", err))
	}

	return config
}
