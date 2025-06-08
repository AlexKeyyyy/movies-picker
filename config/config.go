package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DBUrl           string
	JWTSecret       string
	KinopoiskApiKey string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, reading environment")
	}
	return &Config{
		Port:            os.Getenv("PORT"),
		DBUrl:           os.Getenv("DB_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		KinopoiskApiKey: os.Getenv("KINOPOISK_API_KEY"),
	}
}
