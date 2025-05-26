package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
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

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASS", ""),
		DBName:     getEnv("DB_NAME", "yuk_lomba"),
	}
}
