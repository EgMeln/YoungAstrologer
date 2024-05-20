package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest"
)

var (
	db *sql.DB

	imageRep ImageManager
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("couldn't not connect to docker: %s", err)
	}
	resource, err := pool.Run("postgres", "14.1-alpine", []string{"POSTGRES_PASSWORD=password123"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	var dbHostPort string

	err = pool.Retry(func() error {
		dbHostPort = resource.GetHostPort("5432/tcp")
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:password123@%s/postgres?sslmode=disable", dbHostPort))
		if err != nil {
			return err
		}

		if err := db.Ping(); err != nil {
			return err
		}

		_, err = db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
		if err != nil {
			return err
		}

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return err
		}

		sourceDriver, err := (&file.File{}).Open("file://../../migrations")
		if err != nil {
			return err
		}

		m, err := migrate.NewWithInstance("file", sourceDriver, "postgres", driver)
		if err != nil {
			return err
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Couldn't not connect to db: %s", err)
	}
	defer db.Close()

	imageRep = NewImageManager(db)

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Couldn't purge resource: %s", err)
	}
	err = resource.Expire(1)
	if err != nil {
		log.Fatalf("Couldn't expire: %s", err)
	}

	os.Exit(code)
}
