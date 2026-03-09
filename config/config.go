package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct{
	App AppConfig
	DB 	DBConfig
	JWT JWTConfig
}

type AppConfig struct{
	Name	string
	Port	int
}

type JWTConfig struct{
	Secret 		 string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type DBConfig struct{
	Host 		 	string
	Port 		 	int
	User 		 	string
	Password	string
	Name 		 	string
	Pool 		 	DBPoolConfig
}

type DBPoolConfig struct{
	MaxIdleConnection      int
	MaxOpenConnection      int
	MaxLifetimeConnection  time.Duration
	MaxIdletimeConnection  time.Duration
}

func LoadConfig(envPath ...string) (*Config, error){
	path := ".env"
	if len(envPath) > 0 {
		path = envPath[0]
	}

	_ = godotenv.Load(path)
	
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	accessTTL, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
	}

	refreshTTL, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
	}

	maxIdle, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE: %w", err)
	}

	maxOpen, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN: %w", err)
	}

	maxLifetime, err := time.ParseDuration(os.Getenv("DB_MAX_LIFETIME"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_LIFETIME: %w", err)
	}

	maxIdleTime, err := time.ParseDuration(os.Getenv("DB_MAX_IDLE_TIME"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_TIME: %w", err)
	}

	return &Config{
		App: AppConfig{
			Name: os.Getenv("APP_NAME"),
			Port: port,
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Pool: DBPoolConfig{
				MaxIdleConnection:     maxIdle,
				MaxOpenConnection:     maxOpen,
				MaxLifetimeConnection: maxLifetime,
				MaxIdletimeConnection: maxIdleTime,
			},
		},
		JWT: JWTConfig{
			Secret:     os.Getenv("JWT_SECRET"),
			AccessTTL:  accessTTL,
			RefreshTTL: refreshTTL,
		},
	}, nil
}

