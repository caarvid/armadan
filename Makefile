ifneq (,$(wildcard ./.env))
	include .env
	export
endif

GIT_VERSION ?= $(shell git describe --abbrev=8 --tags --always --dirty)
SERVICE_NAME ?= armadan

export BUILD_VERSION := $(GIT_VERSION)

.PHONY: clean
clean:
	@rm -rf ./tmp ./dist

.PHONE: install
install:
	@go get ./...
	@npm ci

## DEV ##
.PHONY: dev/templ
dev/templ:
	@templ generate --watch --proxy="http://localhost:$(PORT)" --open-browser=false

.PHONY: dev/css
dev/css:
	@npx --yes @tailwindcss/cli -i ./web/css/style.css -o ./web/static/main.css --minify --watch

.PHONY: dev/sql
dev/sql:
	@go run github.com/air-verse/air@v1.61.1 \
		--build.cmd "sqlc generate" \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_dir "db" \
		--build.include_ext "sql" \
		--log.silent "true"

.PHONY: dev/sync
dev/sync:
	@go run github.com/air-verse/air@v1.61.1 \
		--build.cmd "templ generate --notify-proxy" \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_dir "web/static" \
		--build.include_ext "js,css" \
		--log.silent "true"

.PHONY: dev/server
dev/server:
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

.PHONY: dev
dev:
	make -j5 dev/templ dev/css dev/sql dev/sync dev/server 

### BUILD ###
.PHONY: build/css
build/css: 
	@npx --yes @tailwindcss/cli -i ./web/css/style.css -o ./web/static/main.css --minify

.PHONY: build/templ
build/templ:  
	@go run github.com/a-h/templ/cmd/templ@latest generate

.PHONY: build/sql
build/sql: 
	@go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate 

.PHONY: build/armadan
build: clean build/css build/templ build/sql
	@GOOS=linux GOARCH=amd64 go build -o ./dist/armadan ./cmd/armadan/main.go	

.PHONY: build/usertool
build/usertool: 
	@go build -o ./dist/usertool ./cmd/usertool/main.go	

### DOCKER ###
.PHONY: docker/build
docker/build: build
	@docker build \
		--build-arg BUILD_VERSION=$(GIT_VERSION) \
		-f ./Dockerfile \
		-t $(SERVICE_NAME) .
	@docker tag $(SERVICE_NAME):latest $(SERVICE_NAME):$(GIT_VERSION)

### CI/CD ###
.PHONY: ci/build
ci/build: clean install
	@GOOS=linux GOARCH=amd64 go build -o ./dist/armadan ./cmd/armadan/main.go	

### MIGRATIONS ###
.PHONY: migrate/new
migrate/new:
	@goose -dir ./db/migrations sqlite3 $(DB_PATH) create $(name) sql

.PHONY: migrate/up
migrate/up:
	@goose -dir ./db/migrations sqlite3 $(DB_PATH) up

.PHONY: migrate/down
migrate/down:
	@goose -dir ./db/migrations sqlite3 $(DB_PATH) down

