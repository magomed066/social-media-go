package main

import (
	"social-media-app/internal/db"
	"social-media-app/internal/env"
	"social-media-app/internal/handlers"

	"charm.land/log/v2"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {
	initialize()

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:root@localhost:5433/social-media?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENVIRONMENT", "development"),
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Info("db connection pool established")

	handlers := handlers.NewHandlers(db)

	app := &application{
		config:   cfg,
		handlers: handlers,
	}

	if err := app.run(app.mount()); err != nil {
		log.Fatal(err)
	}
}

func initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file")
	}
}
