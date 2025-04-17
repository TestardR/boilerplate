package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"boilerplate/internal/infrastructure/persistence/postgres"

	"github.com/90poe/envconfig"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type postgresSuite struct {
	suite.Suite
	config            postgres.Config
	postgresContainer testcontainers.Container
	db                *sqlx.DB
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(postgresSuite))
}

func (s *postgresSuite) SetupSuite() {
	t := s.T()

	ctx := context.Background()

	err := envconfig.Process("", &s.config)
	require.NoError(t, err)

	err = s.buildPostgresContainer(ctx)
	require.NoError(t, err)

	err = s.buildPostgresConfig(ctx)
	require.NoError(t, err)

	err = s.executeMigrations()
	require.NoError(t, err)

	err = s.createDatabaseHandle()
	require.NoError(t, err)
}

func (s *postgresSuite) buildPostgresContainer(ctx context.Context) error {
	request := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        s.config.Image,
			ExposedPorts: []string{fmt.Sprintf("%s/tcp", s.config.Port)},
			WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			Env: map[string]string{
				"POSTGRES_USER":     s.config.User,
				"POSTGRES_PASSWORD": s.config.Password,
				"POSTGRES_DB":       s.config.DB,
			},
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, request)
	if err != nil {
		return fmt.Errorf("postgres container start failed (did not you forget to run make docker-pull?): %w", err)
	}
	s.postgresContainer = container

	return nil
}

func (s *postgresSuite) buildPostgresConfig(ctx context.Context) error {
	host, err := s.postgresContainer.Host(ctx)
	if err != nil {
		return err
	}

	port, err := s.postgresContainer.MappedPort(ctx, nat.Port(s.config.Port))
	if err != nil {
		return err
	}

	s.config.Host = host
	s.config.Port = port.Port()
	s.config.MigrationsPath = "../../../../db/migrations"

	return nil
}

func (s *postgresSuite) executeMigrations() error {
	migration, err := migrate.New(s.config.SourceURL(), s.config.DBURL())
	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (s *postgresSuite) createDatabaseHandle() error {
	db, err := sqlx.Open("postgres", s.config.DBURL())
	if err != nil {
		return err
	}

	s.db = db

	return nil
}
