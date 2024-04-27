include ./.env
export

init:
	@touch .env

dev:
	@./scripts/clean.sh
	@./scripts/dev.sh

clean:
	@./scripts/clean.sh

gen-template:
	@templ generate

gen-sql:
	@sqlc generate

create-migration:
	@goose -dir ./db/migrations create $(name) sql

migrations-up:
	@goose -dir ./db/migrations postgres "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" up

