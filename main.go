package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

//	@title			ouchi-ipam API
//	@version		0.1
//	@description	ouchi-ipam server
//	@BasePath		/
func main() {
	pgConf := new(PGConfig)
	err := envconfig.Process("", pgConf)
	if err != nil {
		log.Fatal("Error: Failed to read postgres environment variables: ", err)
	}

	pgUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", pgConf.PGUsername, pgConf.PGPassword, pgConf.PGHost, pgConf.PGPort, pgConf.PGDatabase, pgConf.PGSSLMode)
	pool, err := NewPool(pgUrl)
	if err != nil {
		log.Fatal("Error: Failed to connect to database: ", err)
	}
	defer pool.Close()

	serverConf := new(ServerConfig)
	err = envconfig.Process("", serverConf)
	if err != nil {
		log.Fatal("Error: Failed to read web server environment variables: ", err)
	}

	server, err := NewSever(pool)
	if err != nil {
		log.Fatal("Error: Failed to create server: ", err)
	}

	log.Fatal("Error: Server failed to start: ", server.Start(":"+serverConf.Port))
}
