package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            string        `envconfig:"DB_PORT" default:"5432"`
	DB              string        `envconfig:"DB_DB" default:"postgres"`
	User            string        `envconfig:"DB_USER" default:"user"`
	Password        string        `envconfig:"DB_PASSWORD" default:"password"`
	SSL             string        `envconfig:"DB_SSL" default:"disable"`
	Image           string        `envconfig:"POSTGRES_IMAGE" default:"postgres:17-alpine"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"2"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"1"`
	ConnMaxIdleTime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"1m"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_IDLE_TIME" default:"10m"`
	MigrationsPath  string        `envconfig:"DB_MIGRATIONS_PATH" default:"db/migrations"`
}

func (c Config) SourceURL() string {
	return fmt.Sprintf("file://%s", c.MigrationsPath)
}

func (c Config) DBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DB, c.SSL)
}
