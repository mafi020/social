include .env
MIGRATION_PATH = ./cmd/migrate/migrations

.PHONY: install-golang-migrate db-create migration migrate-up migrate-down

install-golang-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

db-create:
	psql -U postgres -h localhost -c "CREATE DATABASE golang_social;"

migration:
	migrate create -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	migrate -database "$(PSQL_URL)" -path $(MIGRATION_PATH) up

migrate-down:
	migrate -database "$(PSQL_URL)" -path $(MIGRATION_PATH) down


