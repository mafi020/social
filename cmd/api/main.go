package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mafi020/social/internal/db"
	"github.com/mafi020/social/internal/env"
	"github.com/mafi020/social/internal/store"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	cfg := config{
		port: env.GetEnvOrPanic("PORT"),
		db: dbConfig{
			url:          env.GetEnvOrPanic("PSQL_URL"),
			maxOpenConns: env.GetEnvAsIntOrPanic("PSQL_MAX_OPEN_CONNS"),
			maxIdleConns: env.GetEnvAsIntOrPanic("PSQL_MAX_IDLE_CONNS"),
			maxIdleTime:  env.GetEnvOrPanic("PSQL_MAX_IDLE_TIME"),
		},
		env: env.GetEnvOrPanic("ENVIRONMENT"),
	}

	db, err := db.New(cfg.db.url, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Printf("Postgres Database Connected")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	log.Fatal(app.start(app.mount()))
}
