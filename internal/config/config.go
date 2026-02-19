package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration loaded from env vars / .env file.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Log      LogConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Port    string
	Env     string
	Version string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type LogConfig struct {
	Level string
}

type RateLimitConfig struct {
	Public        int
	Authenticated int
}

// DSN returns the PostgreSQL Data Source Name.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// Addr returns redis address in host:port format.
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

// Load reads configuration from .env file and environment variables.
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read .env file (ignore error if not found — env vars will be used)
	_ = viper.ReadInConfig()

	cfg := &Config{
		App: AppConfig{
			Port:    viper.GetString("APP_PORT"),
			Env:     viper.GetString("APP_ENV"),
			Version: viper.GetString("APP_VERSION"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
		Log: LogConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
		RateLimit: RateLimitConfig{
			Public:        viper.GetInt("RATE_LIMIT_PUBLIC"),
			Authenticated: viper.GetInt("RATE_LIMIT_AUTHENTICATED"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.App.Port == "" {
		c.App.Port = "8080"
	}
	if c.App.Version == "" {
		c.App.Version = "1.0.0"
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Redis.Host == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}
	return nil
}
