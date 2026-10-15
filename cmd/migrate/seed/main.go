package main

import (
	"errors"
	"log"
	"os"
	"social-media-app/internal/db"
	"social-media-app/internal/env"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	conn, err := db.New(
		env.GetString("DB_ADDR", "postgres://postgres:root@localhost:5433/social-media?sslmode=disable"),
		1, 1, "15m",
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := db.Seed(conn); err != nil {
		return err
	}
	log.Print("Seed complete: 100 users, 200 posts, and 500 comments ensured (existing seed records reused).")
	return nil
}
