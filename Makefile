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