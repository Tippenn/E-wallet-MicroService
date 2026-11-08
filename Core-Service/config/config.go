package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Database *Database
		Server   *Server
		Redis    *Redis
		RabbitMQ *RabbitMQ
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DbName   string
		SslMode  string
		TimeZone string
	}

	Server struct {
		Port    int
		BaseURL string
	}

	Redis struct {
		Host     string
		Port     int
		Password string
	}
	RabbitMQ struct {
		Host              string
		Port              int
		User              string
		Password          string
		VHost             string
		Exchange          string
		PaymentQueue      string
		PaymentRoutingKey string
		StatusQueue       string
		StatusRoutingKey  string
	}
)

var (
	once   sync.Once
	config *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		err := viper.Unmarshal(&config)
		if err != nil {
			panic(err)
		}
	})

	return config
}
