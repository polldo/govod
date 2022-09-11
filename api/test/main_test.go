package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/polldo/govod/config"
	"github.com/polldo/govod/database"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	d, purge, err := startDB()
	if err != nil {
		fmt.Printf("cannot run postgres in a docker container: %s", err)
		os.Exit(1)
	}

	db = d
	code := m.Run()

	if err := purge(); err != nil {
		fmt.Printf("%s", err)
	}

	os.Exit(code)
}

func startDB() (*sqlx.DB, func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_USER=test_user",
			"POSTGRES_PASSWORD=test_pass",
			"POSTGRES_DB=test_db",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		return nil, nil, fmt.Errorf("could not start db container: %w", err)
	}

	container.Expire(60)

	var db *sqlx.DB
	pool.MaxWait = 60 * time.Second
	err = pool.Retry(func() error {
		db, err = database.Open(config.DB{
			User:       "test_user",
			Password:   "test_pass",
			Name:       "test_db",
			Host:       container.GetHostPort("5432/tcp"),
			DisableTLS: true,
		})
		return err
	})

	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to db: %w", err)
	}

	purge := func() error {
		if err := pool.Purge(container); err != nil {
			return fmt.Errorf("cleaning container %s: %w", container.Container.Name, err)
		}
		return nil
	}

	if err := database.Migrate(db); err != nil {
		purge()
		return nil, nil, fmt.Errorf("cannot complete migration on new db: %w", err)
	}

	return db, purge, nil
}
