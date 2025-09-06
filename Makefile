.PHONY: dev up down migrate lint test tidy



# Переменные
ENV_FILE := .env.local
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

# Цвета
COLOR_RESET = \033[0m
COLOR_RED = \033[31m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m

help: ## Показать доступные команды
	@echo "$(COLOR_GREEN)Usage:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/.*Makefile://' | awk 'BEGIN {FS = ":.*?## "}; {printf "$(COLOR_YELLOW)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'


		
lint: ## Run linter
	golangci-lint run

migrate-up: ## Run migrations
	goose -dir ./migrations postgres "${POSTGRES_DSN}" up

migrate-down: ## Run migrations rollback
	goose -dir ./migrations postgres "${POSTGRES_DSN}" down

migrate-status: ## Run migrations status
	goose -dir ./migrations postgres "${POSTGRES_DSN}" status

migrate-create: ## Create new migration
	@echo "Usage: make migrate-create name=<migration_name>"
	@[ "$(name)" ] || (echo "❗  нужно указать name=<migration_name>"; exit 1)
	goose -dir ./migrations create $(name) sql	