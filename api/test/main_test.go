package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/polldo/govod/api"
	"github.com/polldo/govod/config"
	"github.com/polldo/govod/database"
	"github.com/sirupsen/logrus"
)

var dbTest config.DB

func TestMain(m *testing.M) {
	dbTest = config.DB{
		User:       "test_user",
		Password:   "test_pass",
		Name:       "test_db",
		DisableTLS: true,
	}

	c, purge, err := startDB(dbTest.User, dbTest.Password, dbTest.Name)
	if err != nil {
		fmt.Printf("cannot run postgres in a docker container: %s", err)
		os.Exit(1)
	}

	dbTest.Host = c.GetHostPort("5432/tcp")

	code := m.Run()

	if err := purge(); err != nil {
		fmt.Printf("%s", err)
	}

	os.Exit(code)
}

func startDB(user string, pass string, dbname string) (*dockertest.Resource, func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + pass,
			"POSTGRES_DB=" + dbname,
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
	db.Close()

	purge := func() error {
		if err := pool.Purge(container); err != nil {
			return fmt.Errorf("cleaning container %s: %w", container.Container.Name, err)
		}
		return nil
	}

	return container, purge, nil
}

type TestEnv struct {
	*httptest.Server
}

func NewTestEnv(t *testing.T, dbname string) (*TestEnv, error) {

	// Create a new database for each env.
	dbMain, err := database.Open(dbTest)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}

	if _, err := dbMain.ExecContext(context.Background(), "CREATE DATABASE "+dbname); err != nil {
		return nil, fmt.Errorf("creating database %s: %v", dbname, err)
	}
	dbMain.Close()

	// Connect to the new db and perform migrations.
	dbt := dbTest
	dbt.Name = dbname
	dbEnv, err := database.Open(dbt)

	if err := database.Migrate(dbEnv); err != nil {
		return nil, fmt.Errorf("cannot complete migration on new db: %v", err)
	}

	// Redirect log to stdout.
	log := logrus.New()
	log.SetOutput(os.Stdout)

	// Init a new session for authentications.
	sess := scs.New()
	sess.Lifetime = 24 * time.Hour

	api := api.APIMux(api.APIConfig{
		Log:     log,
		DB:      dbEnv,
		Session: sess,
	})

	ts := httptest.NewTLSServer(api)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	ts.Client().Jar = jar
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestEnv{Server: ts}, nil
}
