package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig() *Config {
	viper.SetConfigName("app-config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	return &cfg
}

type Config struct {
	Port int
	DB   DBConfig
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}
