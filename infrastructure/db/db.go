package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rizqdwan/go-mono-project/config"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg config.DBConfig) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.Pool.MaxIdleConnection)
	db.SetMaxOpenConns(cfg.Pool.MaxOpenConnection)
	db.SetConnMaxLifetime(cfg.Pool.MaxLifetimeConnection)
	db.SetConnMaxIdleTime(cfg.Pool.MaxIdletimeConnection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
