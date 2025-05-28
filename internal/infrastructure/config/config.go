package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Server struct {
	Port string
}

type Auth struct {
	JWTSecret string
	*oauth2.Config
}

type Config struct {
	Server
	DB
	Auth
}

var JwtSecret string

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}
	JwtSecret = getEnv("JWT_SECRET", "secret")

	return &Config{
		Server: Server{
			Port: getEnv("APP_PORT", "8080"),
		},
		DB: DB{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", "postgres"),
			Name:     getEnv("DB_NAME", "yuklomba"),
		},
		Auth: Auth{
			JWTSecret: getEnv("JWT_SECRET", "secret"),
			Config: &oauth2.Config{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint: google.Endpoint,
			},
		},
	}, nil
}
