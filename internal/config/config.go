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
		Portal Portal
		Auth   Auth
	}

	Server struct {
		Host           string
		Port           string
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		MaxHeaderBytes int
	}

	Portal struct {
		URL                    string
		AttendanceURL          string
		PerformanceSubjectsURL string
		PerformanceScoreURL    string
	}

	Auth struct {
		ServiceURL string
		Timeout    time.Duration
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

	viper.AutomaticEnv()

	viper.BindEnv("auth.serviceurl", "AUTH_SERVICE_URL")
	viper.BindEnv("portal.attendanceurl", "PORTAL_ATTENDANCE_URL")
	viper.BindEnv("portal.performancesubjectsurl", "PORTAL_PERFORMANCE_SUBJECTS_URL")
	viper.BindEnv("portal.performancescoreurl", "PORTAL_PERFORMANCE_SCORE_URL")

	return viper.ReadInConfig()
}
