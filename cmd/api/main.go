package main

import (
	"social-media-app/internal/env"

	"charm.land/log/v2"
	"github.com/joho/godotenv"
)

func main() {
	initialize()

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	err := app.run(app.mount())
	if err != nil {
		log.Fatal(err)
	}
}

func initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file")
	}
}
