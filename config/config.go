package config

import (
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
	Name 			string
	Port 			int
}

type JWTConfig struct{
	Secret 		 string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type DBConfig struct{
	Host 		 string
	Port 		 int
	User 		 string
	Password string
	Name 		 string
	Pool 		 DBPoolConfig
}

type DBPoolConfig struct{
	MaxIdleConnection      int
	MaxOpenConnection      int
	MaxLifetimeConnection  time.Duration
	MaxIdletimeConnection  time.Duration
}

var Cfg *Config

func LoadConfig(envPath ...string) error{
	
	path := ".env"
	if len(envPath) > 0 {
		path = envPath[0]
	}

	if err := godotenv.Load(path); err != nil {
		return err
	}
	
	port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	accessTTL, _ := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
	refreshTTL, _ := time.ParseDuration(os.Getenv("JWT_REFRESH_TTL"))
	maxIdle, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE"))
	maxOpen, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN"))
	maxLifetime, _ := time.ParseDuration(os.Getenv("DB_MAX_LIFETIME"))
	maxIdleTime, _ := time.ParseDuration(os.Getenv("DB_MAX_IDLE_TIME"))
	
	Cfg = &Config{
		App: AppConfig{
			Name: os.Getenv("APP_NAME"),
			Port: port,
		},
		DB: DBConfig{
			Host: 		os.Getenv("DB_HOST"),
			Port: 		dbPort,
			User: 		os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name: 		os.Getenv("DB_NAME"),
			Pool: DBPoolConfig{
				MaxIdleConnection:      maxIdle,
				MaxOpenConnection:      maxOpen,
				MaxLifetimeConnection:  maxLifetime,
				MaxIdletimeConnection:  maxIdleTime,
			},
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			AccessTTL: accessTTL,
			RefreshTTL: refreshTTL,
		},
	}

	return nil
}

