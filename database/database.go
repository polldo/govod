// This file has been adapted from the excellent ardanlabs repo:
// https://github.com/ardanlabs/service ardanlabs .
package database

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/polldo/govod/config"
)

// lib/pq errorCodeNames
// https://github.com/lib/pq/blob/master/error.go#L178
const uniqueViolation = "23505"

// Set of error variables for CRUD operations.
var (
	ErrDBNotFound        = errors.New("not found")
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
)

//go:embed sql/migration/*.sql
var sqlFS embed.FS

func Migrate(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	fs, err := iofs.New(sqlFS, "sql/migration")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", fs, "postgres", driver)
	if err != nil {
		return err
	}

	return m.Up()
}

func Seed(db *sqlx.DB, seed string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seed); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func Transaction(db *sqlx.DB, f func(db sqlx.ExtContext) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot begin the transaction: %w", err)
	}

	if err := f(tx); err != nil {
		if terr := tx.Rollback(); terr != nil {
			return fmt.Errorf("transaction failed but could not rollback: %v: %w", terr, err)
		}
		return fmt.Errorf("transaction failed and rolled back: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit the transaction: %w", err)
	}

	return nil
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg config.DB) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NamedExecContext is a helper function to execute a CUD operation.
func NamedExecContext(ctx context.Context, db sqlx.ExtContext, query string, data any) error {
	if _, err := sqlx.NamedExecContext(ctx, db, query, data); err != nil {

		// Checks if the error is of code 23505 (unique_violation).
		if pqerr, ok := err.(*pq.Error); ok && pqerr.Code == uniqueViolation {
			return ErrDBDuplicatedEntry
		}
		return err
	}

	return nil
}

// NamedQuerySlice is a helper function for executing queries that return a
// collection of data to be unmarshalled into a slice.
func NamedQuerySlice[T any](ctx context.Context, db sqlx.ExtContext, query string, data any, dest *[]T) error {
	rows, err := sqlx.NamedQueryContext(ctx, db, query, data)
	if err != nil {
		return err
	}
	defer rows.Close()

	var slice []T
	for rows.Next() {
		v := new(T)
		if err := rows.StructScan(v); err != nil {
			return err
		}
		slice = append(slice, *v)
	}

	if slice != nil {
		*dest = slice
	}

	return nil
}

// NamedQueryStruct is a helper function for executing queries that return a
// single value to be unmarshalled into a struct type.
func NamedQueryStruct(ctx context.Context, db sqlx.ExtContext, query string, data any, dest any) error {
	rows, err := sqlx.NamedQueryContext(ctx, db, query, data)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrDBNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}
