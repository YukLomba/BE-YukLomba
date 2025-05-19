package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort    int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func LoadConfig() *Config {
	appPort, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))

	return &Config{
		AppPort:    appPort,
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASS", ""),
		DBName:     getEnv("DB_NAME", "go_fiber_clean_arch"),
	}
}
