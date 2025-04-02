ifneq (,$(wildcard ./.env))
	include .env
	export
endif

GIT_VERSION ?= $(shell git describe --abbrev=8 --tags --always --dirty)
SERVICE_NAME ?= armadan

export BUILD_VERSION := $(GIT_VERSION)

.PHONY: clean
clean:
	@go clean
	@rm -rf ./tmp ./dist

.PHONY: install/hooks
install/hooks:
	@echo "Installing Git hooks..."
	@cp scripts/pre-commit .git/hooks/pre-commit
	@cp scripts/pre-push .git/hooks/pre-push
	@chmod +x .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-push
	@echo "Git hooks installed."

.PHONY: install/deps
install/deps:
	@go get ./...
	@pnpm install --frozen-lockfile

.PHONY: setup
setup: clean install/hooks install/deps
	@./scripts/setup

.PHONY: release
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is not set."; \
		echo "Usage: make release VERSION=x.y.z"; \
		exit 1; \
	fi
	git tag -a "v$(VERSION)" -m "Release v$(VERSION)"
	git push origin "v$(VERSION)"

## DEV ##
.PHONY: dev/templ
dev/templ:
	@templ generate --watch

.PHONY: dev/css
dev/css:
	@npx @tailwindcss/cli -i ./web/css/style.css -o ./web/static/main.css --minify --watch

.PHONY: dev/sql
dev/sql:
	@air \
		--build.cmd "sqlc generate" \
		--build.bin "true" \
		--build.delay "100" \
		--build.include_dir "db" \
		--build.include_ext "sql" \
		--log.silent "true"

.PHONY: dev/sync
dev/sync:
	@air \
		--build.cmd "templ generate --notify-proxy" \
		--build.bin "true" \
		--build.delay "100" \
		--build.include_dir "web/static" \
		--build.include_ext "js,css" \
		--log.silent "true"

.PHONY: dev/server
dev/server:
	@air \
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
	make -j4 dev/templ dev/css dev/sql dev/server 

### TESTING ###
.PHONY: test
test:
	@go test ./...	

### BUILD ###
.PHONY: build/css
build/css: 
	@npx @tailwindcss/cli -i ./web/css/style.css -o ./web/static/main.css --minify

.PHONY: build/html
build/html:  
	@templ generate

.PHONY: build/sql
build/sql: 
	@sqlc generate 

.PHONY: build
build: clean build/css build/html build/sql 
	@go build -o dist/armadan ./cmd/armadan/main.go

### TOOLS ###
.PHONY: tools/create_user
tools/create_user:
	@go run ./tools/create_user.go --email=$(EMAIL) --password=$(PASSWORD) --role=$(ROLE) 

### MIGRATIONS ###
.PHONY: migrate/new
migrate/new:
	@goose -dir ./db/migrations sqlite3 $(DB_PATH_MIGRATIONS) create $(name) sql

.PHONY: migrate/up
migrate/up:
	@goose -dir ./db/migrations sqlite3 $(DB_PATH_MIGRATIONS) up

.PHONY: migrate/down
migrate/down:
	@goose -dir ./db/migrations sqlite3 $(DB_PATH_MIGRATIONS) down

## HOOKS ##
.PHONY: hooks/pre-commit
hooks/pre-commit: build/css build/sql build/html
	@git add ./web/static/main.css
