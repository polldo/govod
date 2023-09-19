package test

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"math/rand"
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
	"github.com/plutov/paypal/v4"
	"github.com/polldo/govod/api"
	"github.com/polldo/govod/api/background"
	"github.com/polldo/govod/config"
	"github.com/polldo/govod/database"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v74"
	stripecl "github.com/stripe/stripe-go/v74/client"
	"golang.org/x/crypto/bcrypt"
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

	// rand is used to generate random IDs and names in tests.
	// Let's use a different random seed for each test execution.
	rand.Seed(time.Now().UTC().UnixNano())

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

	container.Expire(120)

	var db *sqlx.DB
	pool.MaxWait = 120 * time.Second
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

type mockMailer struct {
	token string
}

func (m *mockMailer) SendActivationToken(token string, dst string) error {
	m.token = token
	return nil
}

func (m *mockMailer) SendResetToken(token string, dst string) error {
	m.token = token
	return nil
}

const seedTest = `
INSERT INTO users (user_id, name, email, role, active, password_hash, created_at, updated_at) VALUES
	('ae127240-ce13-4789-aafd-d2f31e7ee487', 'Admin', '{{ .AdminEmail}}', 'ADMIN', TRUE, '{{ .AdminPassHash}}', '2022-09-16 00:00:00', '2022-09-16 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Test', '{{ .UserEmail}}', 'USER', TRUE, '{{ .UserPassHash}}', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;
`

type TestEnv struct {
	*httptest.Server

	// Admin test credentials in seed.
	AdminEmail string
	AdminPass  string

	// User test credentials in seed.
	UserEmail string
	UserPass  string

	// Collect mocked dependencies here to make them
	// available to all tests.
	Mailer        *mockMailer
	Paypal        *mockPaypal
	Stripe        *mockStripe
	WebhookSecret string
}

func (te *TestEnv) parseSeed() (string, error) {
	tmp := struct {
		AdminEmail    string
		AdminPassHash string
		UserEmail     string
		UserPassHash  string
	}{
		AdminEmail: te.AdminEmail,
		UserEmail:  te.UserEmail,
	}

	h, err := bcrypt.GenerateFromPassword([]byte(te.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	tmp.AdminPassHash = string(h)

	h, err = bcrypt.GenerateFromPassword([]byte(te.UserPass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	tmp.UserPassHash = string(h)

	t, err := template.New("seed").Parse(seedTest)
	if err != nil {
		return "", err
	}

	var res bytes.Buffer
	if err = t.Execute(&res, tmp); err != nil {
		return "", err
	}

	return res.String(), nil
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

	// Setup the test environment with some users.
	te := &TestEnv{
		AdminEmail: "admin@govod.com",
		AdminPass:  "admin-password123",
		UserEmail:  "user@govod.com",
		UserPass:   "user-password123",
	}

	seed, err := te.parseSeed()
	if err != nil {
		return nil, err
	}

	// Apply the seed to the new db.
	if err := database.Seed(dbEnv, seed); err != nil {
		return nil, fmt.Errorf("cannot init db with seed: %v", err)
	}

	// Redirect log to stdout.
	log := logrus.New()
	log.SetOutput(os.Stdout)

	// Init a new session for authentications.
	sess := scs.New()
	sess.Lifetime = 24 * time.Hour

	// Build a mocked mailer to allow signup in tests.
	mail := &mockMailer{}
	te.Mailer = mail

	// Init a background manager to safely spawn go-routines.
	bg := background.New(log)

	// Setup the mock for paypal payments.
	te.Paypal = &mockPaypal{}
	ppserver := httptest.NewServer(te.Paypal.handle())

	// Build the paypal client to allow payments and make it pointing to the mocked server.
	// No need to generate a token since the server is mocked.
	pp, err := paypal.NewClient("test", "test", ppserver.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to build the paypal client: %w", err)
	}

	// Setup the mock for stripe payments.
	te.Stripe = &mockStripe{}
	strpserver := httptest.NewServer(te.Stripe.handle())

	// Build the stripe client to allow payments.
	strpcfg := config.Stripe{
		APISecret:     "random-key",
		WebhookSecret: "random-test-secret",
		SuccessURL:    "/success.html",
		CancelURL:     "/cart.html",
	}
	te.WebhookSecret = strpcfg.WebhookSecret
	strp := &stripecl.API{}

	// Point to the mocked stripe server.
	strp.Init(strpcfg.APISecret, &stripe.Backends{
		API:     stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: &strpserver.URL}),
		Connect: stripe.GetBackend(stripe.ConnectBackend),
		Uploads: stripe.GetBackend(stripe.UploadsBackend),
	})

	api := api.APIMux(api.APIConfig{
		CorsOrigin: "",
		Log:        log,
		DB:         dbEnv,
		Session:    sess,
		Mailer:     mail,
		Background: bg,
		Paypal:     pp,
		Stripe:     strp,
		StripeCfg:  strpcfg,
	})

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	te.Server = httptest.NewTLSServer(api)
	te.Server.Client().Jar = jar
	te.Server.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return te, nil
}
