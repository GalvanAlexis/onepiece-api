BINARY_NAME=onepiece-api
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/api/main.go
MIGRATE_PATH=migrations
DB_URL=postgres://onepiece:onepiece_secret@localhost:5432/onepiece_db?sslmode=disable

.PHONY: all build run clean docker-up docker-down migrate-up migrate-down sqlc-gen lint test

## Build & Run
all: build

build:
	@echo ">> Building..."
	go build -o $(BINARY_PATH) $(MAIN_PATH)

run:
	@echo ">> Running..."
	go run $(MAIN_PATH)

clean:
	@echo ">> Cleaning..."
	rm -rf bin/

## Docker
docker-up:
	@echo ">> Starting services..."
	docker-compose up -d postgres redis
	@echo ">> Waiting for postgres to be ready..."
	@sleep 3

docker-down:
	@echo ">> Stopping services..."
	docker-compose down

docker-all:
	@echo ">> Starting all services..."
	docker-compose up -d

## Migrations
migrate-up:
	@echo ">> Running migrations UP..."
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	@echo ">> Running migrations DOWN..."
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down

migrate-create:
	@echo ">> Creating migration: $(name)..."
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

## SQLC
sqlc-gen:
	@echo ">> Generating SQLC code..."
	sqlc generate

## Testing & Linting
test:
	@echo ">> Running tests..."
	go test ./... -v -race

lint:
	@echo ">> Linting..."
	golangci-lint run ./...

## Deps
tidy:
	@echo ">> Running go mod tidy..."
	go mod tidy
