include ./.env
export

init:
	@touch .env

dev: 
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

migrations-down:
	@goose -dir ./db/migrations postgres "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" down

start-db:
	docker compose --file ./docker/db.yml --env-file ./.env up --detach

stop-db:
	docker compose --file ./docker/db.yml --env-file ./.env down 
