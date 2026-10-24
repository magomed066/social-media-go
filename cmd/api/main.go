package main

import (
	"social-media-app/internal/db"
	"social-media-app/internal/env"
	"social-media-app/internal/handlers"
	"social-media-app/internal/mailer"
	"time"

	"charm.land/log/v2"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

//	@title			Social Media API
//	@version		1.0
//	@description	API for social network.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	initialize()

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:3000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:root@localhost:5433/social-media?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENVIRONMENT", "development"),
		mail: mailConfig{
			host:      env.GetString("MAIL_HOST", ""),
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
		},
		frontendUrl: env.GetString("FRONTEND_URL", "localhost:5173"),
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

	mailer, err := mailer.NewWnessenMailer(
		cfg.mail.host,
		cfg.mail.fromEmail,
		env.GetString("SMTP_USERNAME", ""),
		env.GetString("SMTP_PASSWORD", ""),
	)
	if err != nil {
		log.Fatal(err)
	}

	handlers := handlers.NewHandlers(db, mailer, cfg.frontendUrl)

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
