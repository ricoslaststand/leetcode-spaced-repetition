package main

import (
	"flag"

	"leetcode-spaced-repetition/db/migrations"
	"leetcode-spaced-repetition/internal"

	goose "github.com/pressly/goose/v3"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up or down")
	flag.Parse()

	goose.SetBaseFS(migrations.MigrationFiles)
	config, err := internal.GetConfig()
	if err != nil {
		panic(err)
	}
	db, err := internal.GetDBConnFromConfig(&config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = db.Close() }()

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	switch *direction {
	case "up":
		if err := goose.Up(db, "."); err != nil {
			panic(err)
		}
	case "down":
		if err := goose.Down(db, "."); err != nil {
			panic(err)
		}
	default:
		panic("invalid direction: " + *direction + " (must be 'up' or 'down')")
	}
}
