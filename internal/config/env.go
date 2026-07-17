package config

import (
	"fmt"
	"os"
)

type EnvConfig struct {
	AdminUsername string
	AdminPassword string
	Port          string
	DbPath        string
}

func Load() (*EnvConfig, error) {
	cfg := &EnvConfig{
		AdminUsername: os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		Port:          getEnv("PORT", "8080"),
		DbPath:        getEnv("DB_PATH", "./data/nanoshare.db"),
	}

	if cfg.AdminUsername == "" || cfg.AdminPassword == "" {
		return nil, fmt.Errorf("ADMIN_USERNAME and ADMIN_PASSWORD must be set")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
