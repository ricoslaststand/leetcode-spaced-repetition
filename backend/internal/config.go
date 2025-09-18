package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug      bool
	PostgresDB struct {
		Username string `envconfig:"POSTGRES_USERNAME" required:"true"`
		Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
		DB       string `envconfig:"POSTGRES_DB" required:"true"`
	}
}

func GetConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envVars := os.Environ()

	fmt.Printf("%+v\n", envVars)

	var config Config
	err = envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("The config is %+v\n", config)
	// format := "Debug: %v\nPort: %d\nUser: %s\nRate: %f\nTimeout: %s\n"
	// _, err = fmt.Printf(format, s.Debug, s.Port, s.User, s.Rate, s.Timeout)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config, nil
}
