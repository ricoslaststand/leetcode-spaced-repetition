package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetDBConnFromConfig(config Config) (*sql.DB, error) {
	// TODO: Make this configurable for host and port
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"db",
		5432,
		config.PostgresDB.Username, config.PostgresDB.Password, config.PostgresDB.DB)
	return sql.Open("postgres", connStr)
}

func GetDBConnFromDSN(dsn string) (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}
