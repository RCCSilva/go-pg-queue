package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func ConnectDatabase(config Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.ConnectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
