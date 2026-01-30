package db

import (
	"database/sql"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rizqdwan/go-mono-project/config"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg config.DBConfig) (*Database, error) {
	dsn := "postgres://" +
		cfg.User + ":" +
		cfg.Password + "@" +
		cfg.Host + ":" +
		strconv.Itoa(cfg.Port) + "/" +
		cfg.Name

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.Pool.MaxIdleConnection)
	db.SetMaxOpenConns(cfg.Pool.MaxOpenConnection)
	db.SetConnMaxLifetime(cfg.Pool.MaxLifetimeConnection)
	db.SetConnMaxIdleTime(cfg.Pool.MaxIdletimeConnection)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
