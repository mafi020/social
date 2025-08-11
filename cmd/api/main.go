package main

import (
	"github.com/mafi020/social/internal/db"
	"github.com/mafi020/social/internal/env"
	"github.com/mafi020/social/internal/logger"
	"github.com/mafi020/social/internal/store"
)

func main() {
	cfg := &config{
		port: env.GetEnvOrPanic("PORT"),
		db: &dbConfig{
			url:          env.GetEnvOrPanic("PSQL_URL"),
			maxOpenConns: env.GetEnvAsIntOrPanic("PSQL_MAX_OPEN_CONNS"),
			maxIdleConns: env.GetEnvAsIntOrPanic("PSQL_MAX_IDLE_CONNS"),
			maxIdleTime:  env.GetEnvOrPanic("PSQL_MAX_IDLE_TIME"),
		},
		env: env.GetEnvOrPanic("ENVIRONMENT"),
	}

	// Logger: https://github.com/uber-go/zap
	logger := logger.New()
	defer logger.Sync()

	// Configure the Postgres DB
	db, err := db.New(cfg.db.url, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Panicw("Failed to connect to Postgres DB", err)
	}
	defer db.Close()
	logger.Infow("Postgres Database Connected")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	logger.Fatal(app.start(app.mount()))
}
