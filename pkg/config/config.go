package config

import (
	"os"
)

type EnvironmentType string

const (
	Production  EnvironmentType = "prod"
	Development EnvironmentType = "dev"
	Local       EnvironmentType = "local"
)

type Config struct {
	APP         string
	Environment EnvironmentType
	LogLevel    string

	Server struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	Bot struct {
		Token      string
		WebhookURL string
	}

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		Sslmode  string
	}

	Redis struct {
		Host            string
		Port            string
		Password        string
		DB              string
		StorageDeadline string
	}

	LLM struct {
		BaseURL   string
		ModelName string
		SecretKey string
		Timeout   string
	}
}

func New() *Config {
	var config Config

	config.APP = getEnv("APP", "dastyor-bot")
	config.Environment = EnvironmentType(getEnv("ENVIRONMENT", "develop"))
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "q")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnv("SERVER_PORT", ":8000")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// bot init
	config.Bot.Token = getEnv("BOT_TOKEN", "")
	config.Bot.WebhookURL = getEnv("BOT_WEBHOOK_URL", "")

	// initialization db
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "dastyor")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "postgres")
	config.DB.Sslmode = getEnv("POSTGRES_SSLMODE", "disable")

	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = getEnv("REDIS_DB", "0")
	config.Redis.StorageDeadline = getEnv("REDIS_STORAGE_DEADLINE", "30m")

	config.LLM.BaseURL = getEnv("LLM_BASE_URL", "")
	config.LLM.ModelName = getEnv("LLM_MODEL_NAME", "")
	config.LLM.SecretKey = getEnv("LLM_SECRET_KEY", "")
	config.LLM.Timeout = getEnv("LLM_TIMEOUT", "5m")

	return &config
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
