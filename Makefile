calc:
	@go build -o calculator/bin/calculator calculator/main.go
	@./calculator/bin/calculator

docker:
	@docker compose -f calculator/docker-compose.yml up -d

clean:
	@docker compose -f calculator/docker-compose.yml down --volumes