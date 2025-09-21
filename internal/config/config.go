package config

import (
	"fmt"
	"time"

	"github.com/anton1ks96/college-app-core/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server Server
	}

	Server struct {
		Host           string
		Port           string
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		MaxHeaderBytes int
	}
)

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found, using system environment variables")
	}

	if err := parseConfigFile("./configs"); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return &cfg, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")
	viper.SetConfigType("yml")

	return viper.ReadInConfig()
}
