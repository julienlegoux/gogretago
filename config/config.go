package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	JWTExpiresIn   string
	ResendAPIKey   string
	ResendFromEmail string
	Port           int
	AppEnv         string
}

var cfg *Config

func Load() (*Config, error) {
	_ = godotenv.Load()

	port, _ := strconv.Atoi(getEnv("PORT", "3000"))

	cfg = &Config{
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiresIn:   getEnv("JWT_EXPIRES_IN", "24h"),
		ResendAPIKey:   getEnv("RESEND_API_KEY", ""),
		ResendFromEmail: getEnv("RESEND_FROM_EMAIL", ""),
		Port:           port,
		AppEnv:         getEnv("APP_ENV", "development"),
	}

	return cfg, nil
}

func Get() *Config {
	if cfg == nil {
		Load()
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
