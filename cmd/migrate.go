package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDatabase(config Config) error {
	m, err := migrate.New("file://db/migrations", config.ConnectionString)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()

	if err == migrate.ErrNoChange {
		return nil
	}

	return err
}
