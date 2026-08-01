package internal

import (
	"log"

	"github.com/google/uuid"
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
	// OwnerUsername is the username the authenticating reverse proxy reports for the
	// single account this deployment belongs to.
	OwnerUsername string `envconfig:"OWNER_USERNAME" required:"true"`
	AppEnv        string `envconfig:"APP_ENV"`
	// OwnerUserID is the user_id every row in the database is written under. This app is
	// single-user by design; see internal.OwnerOnlyAuthMiddleware.
	OwnerUserID uuid.UUID `envconfig:"OWNER_USER_ID" required:"true"`
	AppPort     uint      `envconfig:"APP_PORT"`
	Debug       bool
}

// IsDevelopment reports whether the app is running in a local development environment.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
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
