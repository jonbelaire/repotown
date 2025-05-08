.PHONY: build run test clean migrate-up migrate-down lint

# Building
build:
	go build -o bin/parking-service ./cmd/server

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
	docker build -t parking-service .

docker-run:
	docker run -p 8082:8082 --env-file .env parking-service