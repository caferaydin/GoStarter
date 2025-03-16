package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DB_CONN"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
	return cfg
}

func ConnectDB(connStr string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln("Failed to connect to DB:", err)
	}
	return db
}
