package config

import (
	"os"
	"strconv"
	"time"
)

const defaultJWTExpiryHours = 24

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret     string
	JWTExpiryTime time.Duration
}

func Load() Config {
	return Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSLMODE"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiryTime: jwtExpiry(),
	}
}

func jwtExpiry() time.Duration {
	hours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || hours < 1 {
		hours = defaultJWTExpiryHours
	}
	return time.Duration(hours) * time.Hour
}
