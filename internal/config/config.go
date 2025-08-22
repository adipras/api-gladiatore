package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	DatabaseHost    string
	DatabasePort    int
	DatabaseUser    string
	DatabasePass    string
	DatabaseName    string
	GeneratedPath   string
	TemplatesPath   string
	AllowCORS       bool
}

func Load() (*Config, error) {
	godotenv.Load()

	dbPort, _ := strconv.Atoi(getEnv("DATABASE_PORT", "3306"))

	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		DatabaseHost:    getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:    dbPort,
		DatabaseUser:    getEnv("DATABASE_USER", "root"),
		DatabasePass:    getEnv("DATABASE_PASS", ""),
		DatabaseName:    getEnv("DATABASE_NAME", "api_gladiatore"),
		GeneratedPath:   getEnv("GENERATED_PATH", "./generated"),
		TemplatesPath:   getEnv("TEMPLATES_PATH", "./internal/templates"),
		AllowCORS:       getEnv("ALLOW_CORS", "true") == "true",
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}