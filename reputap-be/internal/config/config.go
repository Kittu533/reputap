package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	AppEnv           string
	DBHost           string
	DBPort           string
	DBName           string
	DBUser           string
	DBPassword       string
	DBSSLMode        string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	AppBaseURL       string
	RedirectBaseURL  string
	UploadDir        string
	MaxUploadSize    int64
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		accessExpiry = 15 * time.Minute
	}

	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		refreshExpiry = 168 * time.Hour
	}

	maxUpload, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	if err != nil {
		maxUpload = 10485760
	}

	return &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		AppEnv:           getEnv("APP_ENV", "development"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5433"),
		DBName:           getEnv("DB_NAME", "reputap"),
		DBUser:           getEnv("DB_USER", "reputap"),
		DBPassword:       getEnv("DB_PASSWORD", "secret"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTAccessExpiry:  accessExpiry,
		JWTRefreshExpiry: refreshExpiry,
		AppBaseURL:       getEnv("APP_BASE_URL", "http://localhost:8080"),
		RedirectBaseURL:  getEnv("REDIRECT_BASE_URL", "http://localhost:8080/r"),
		UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSize:    maxUpload,
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser + " password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=" + c.DBSSLMode
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
