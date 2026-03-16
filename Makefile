.PHONY: run-dev run-api run-worker migrate-up migrate-down docker-up docker-down docker-build

run-dev: docker-up migrate-up
	$(MAKE) run-api & $(MAKE) run-worker

run-api:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/payments?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/payments?sslmode=disable" down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build