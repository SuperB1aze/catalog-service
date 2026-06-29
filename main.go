package main

import (
	"log"

	"github.com/SuperB1aze/catalog-service/internal/app/config"
)

func main() {
	config.Load(&config.Root)

	cfg := config.Root

	log.Printf("Server will start on port: %d", cfg.Processor.WebServer.ListenPort)
	log.Printf("Database: %s@%s/%s",
		cfg.Repository.Postgres.Username,
		cfg.Repository.Postgres.Address,
		cfg.Repository.Postgres.Name)
	log.Printf("Environment: %s, LogLevel: %s",
		cfg.Monitor.Environment,
		cfg.Monitor.LogLevel)
}
