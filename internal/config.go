package internal

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB struct {
		URL      string `envconfig:"POSTGRES_URL" required:"true"`
		Username string `envconfig:"POSTGRES_USERNAME" required:"true"`
		Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
		DB       string `envconfig:"POSTGRES_DB" required:"true"`
	}
	AppEnv  string `envconfig:"APP_ENV"`
	AppPort uint   `envconfig:"APP_PORT"`
	Debug   bool
}

func GetConfig() (Config, error) {
	// .env is optional — in Docker, env vars are injected via --env-file
	_ = godotenv.Load()

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config, nil
}
