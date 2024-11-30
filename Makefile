ifneq (,$(wildcard ./.env))
	include .env
	export
endif

cmd-exists-%:
	@hash $(*) > /dev/null 2>&1 || (echo "ERROR: '$(*)' must be installed and available in your $PATH"; exit 1)

dev-templ:
	@templ generate --watch --proxy="http://localhost:$(PORT)" --open-browser=false

dev-css:
	@npx tailwindcss -i ./web/css/style.css -o ./web/static/main.css --minify --watch

dev-sql:
	@go run github.com/air-verse/air@v1.61.1 \
		--build.cmd "sqlc generate" \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_dir "db" \
		--build.include_ext "sql" \
		--log.silent "true"

dev-sync:
	@go run github.com/air-verse/air@v1.61.1 \
		--build.cmd "templ generate --notify-proxy" \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_dir "web/static" \
		--build.include_ext "js,css" \
		--log.silent "true"

dev-server:
	@go run github.com/air-verse/air@v1.61.1 \
		--build.cmd "go build -o ./tmp/bin/main ./cmd/armadan/main.go" \
		--build.bin "./tmp/bin/main" \
		--build.delay "100" \
		--build.exclude_dir "node_modules,tmp,assets,db,docker,web/css,web/static" \
		--build.exclude_regex "_test.go" \
		--build.send_interrupt "true" \
		--build.include_ext "go" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit "true" \
		--log.silent "true"

dev:
	make -j5 dev-templ dev-css dev-sql dev-sync dev-server 

css: 
	@npx tailwindcss -i ./web/css/style.css -o ./web/static/main.css --minify

gen-templ: cmd-exists-templ
	@templ generate

gen-sql: cmd-exists-sqlc
	@sqlc generate

build-app:
	@GOOS=linux GOARCH=amd64 go build -o ./dist/armadan ./cmd/armadan/main.go	

build: clean css gen-templ gen-sql build-app

new-migration: cmd-exists-goose
	@goose -dir ./db/migrations create $(name) sql

migrate-up: cmd-exists-goose
	@goose -dir ./db/migrations postgres "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" up

migrate-down: cmd-exists-goose
	@goose -dir ./db/migrations postgres "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" down

db-start: cmd-exists-docker
	docker compose --file ./docker/db.yml --env-file ./.env up --detach

db-stop: cmd-exists-docker
	docker compose --file ./docker/db.yml --env-file ./.env down 

clean:
	@rm -rf ./tmp ./dist

