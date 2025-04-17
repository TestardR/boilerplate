package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

type Client struct {
	db *sqlx.DB
}

func NewClient(config Config) (Client, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%t",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DB,
		config.SSL,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return Client{}, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return Client{db: db}, nil
}

func (c Client) DB() *sqlx.DB {
	return c.db
}

func (c Client) HealthCheck(_ context.Context) error {
	return c.db.Ping()
}

func (c Client) Close() error {
	return c.db.Close()
}
