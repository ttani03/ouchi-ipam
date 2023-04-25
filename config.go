package main

type PGConfig struct {
	PGUsername string `envconfig:"PG_USERNAME"`
	PGPassword string `envconfig:"PG_PASSWORD"`
	PGHost     string `envconfig:"PG_HOST"`
	PGPort     string `envconfig:"PG_PORT" default:"5432"`
	PGDatabase string `envconfig:"PG_DATABASE"`
	PGSSLMode  string `envconfig:"PG_SSL_MODE" default:"disable"`
}

type ServerConfig struct {
	Port string `envconfig:"PORT" default:"8080"`
}
