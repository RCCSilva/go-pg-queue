package main

type Config struct {
	ConnectionString string
}

func NewConfig() Config {
	return Config{
		ConnectionString: "postgresql://user:password@localhost:5432/pg_queue?sslmode=disable",
	}
}
