CONFIG_FILE ?= ./config/config.yml
APP_DSN ?= $(shell sed -n 's/^dsn:[[:space:]]*"\(.*\)"/\1/p' $(CONFIG_FILE))
MIGRATE := docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate:v4.14.1 -path=/migrations/ -database "$(APP_DSN)"
DB_CONTAINER ?= ad-service_db_1

.PHONY: default
default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## run the API server
	go run ./cmd/ad-service/main.go

.PHONY: build
build:  ## build the API server binary
	go build -a -o server ./cmd/ad-service

.PHONY: build-docker
build-docker: ## build the API server as a docker image
	docker build -f ./Dockerfile -t server .

.PHONY: migrate
migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1

.PHONY: testdata
testdata: ## add test data
	@echo "Add test data..."
	@docker cp ./testdata/testdata.sql $(DB_CONTAINER):/testdata.sql
	@docker exec -it $(DB_CONTAINER) psql "$(APP_DSN)" -f testdata.sql
