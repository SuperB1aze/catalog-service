package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/SuperB1aze/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

// Root — глобальный доступ к конфигурации
var Root Config

// Load загружает конфигурацию из .env файла и переменных окружения
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal("Error processing env config: ", err)
	}
}
