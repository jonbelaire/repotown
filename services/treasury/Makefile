.PHONY: build run test clean migrate-up migrate-down lint

# Building
build:
	go build -o bin/treasury-service ./cmd/server

# Running
run:
	go run ./cmd/server

# Testing
test:
	go test -v ./...

# Database migrations
migrate-up:
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down:
	migrate -path migrations -database "${DATABASE_URL}" down

# Clean build artifacts
clean:
	rm -rf bin/

# Linting
lint:
	golangci-lint run

# Docker commands
docker-build:
	docker build -t treasury-service .

docker-run:
	docker run -p 8083:8083 --env-file .env treasury-service