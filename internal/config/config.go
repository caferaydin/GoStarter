package config

import (
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          []byte
	RefreshSecret      []byte
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	accessDur, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	refreshDur, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRY"))

	cfg := &Config{
		DatabaseURL:        os.Getenv("DB_CONN"),
		JWTSecret:          []byte(os.Getenv("JWT_SECRET")),
		RefreshSecret:      []byte(os.Getenv("REFRESH_SECRET")),
		AccessTokenExpiry:  accessDur,
		RefreshTokenExpiry: refreshDur,
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
