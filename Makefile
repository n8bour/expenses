.PHONY: help clean

.DEFAULT_GOAL := help

calc: ## Build and run calculator service
	@go build -o calculator/bin/calculator calculator/main.go
	@./calculator/bin/calculator

docker: ## Run the docker containers
	@docker compose -f calculator/docker-compose.yml up -d

clean: ## Remove all docker containers and volumes
	@docker compose -f calculator/docker-compose.yml down --volumes

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
