package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetDBConnFromConfig(config Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.PostgresDB.Username, config.PostgresDB.Password, config.PostgresDB.DB)
	return sql.Open("postgres", connStr)
}
