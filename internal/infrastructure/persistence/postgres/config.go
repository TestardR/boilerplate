package postgres

import "time"

type Config struct {
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            int           `envconfig:"DB_PORT" default:"5432"`
	DB              string        `envconfig:"DB_DB" default:"postgres"`
	User            string        `envconfig:"DB_USER" default:"test"`
	Password        string        `envconfig:"DB_PASSWORD" default:"password"`
	SSL             bool          `envconfig:"DB_SSL" default:"false"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"2"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"1"`
	ConnMaxIdleTime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"1m"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_IDLE_TIME" default:"10m"`
}
