package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Server   ServerConfig
	Security SecurityConfig
	Logging  LoggingConfig
	Upload   UploadConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type ServerConfig struct {
	Port string
	Host string
	Mode string
}

type SecurityConfig struct {
	CORSOrigin        string
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

type LoggingConfig struct {
	Level  string
	Format string
}

type UploadConfig struct {
	MaxFileSize int64
	UploadPath  string
}

func Load() (*Config, error) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("PSQL_HOST", "localhost"),
			Port:     getEnv("PSQL_PORT", "5437"),
			Username: getEnv("PSQL_USER", "postgres"),
			Password: getEnv("PSQL_PASSWORD", "password"),
			DBName:   getEnv("PSQL_DBNAME", "onlinechat"),
			SSLMode:  getEnv("PSQL_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Address:  getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret-change-this"),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ":8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Security: SecurityConfig{
			CORSOrigin:        getEnv("CORS_ORIGIN", "*"),
			RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", "1m"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Upload: UploadConfig{
			MaxFileSize: getEnvAsInt64("MAX_FILE_SIZE", 10485760), // 10MB
			UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
		},
	}

	if config.JWT.Secret == "default-secret-change-this" {
		return nil, fmt.Errorf("JWT_SECRET must be set to a secure value")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Minute
}
