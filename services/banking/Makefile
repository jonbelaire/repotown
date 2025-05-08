.PHONY: all build run test clean docker

# Binary name
BINARY_NAME=banking-service

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get

# Main package
MAIN_PACKAGE=./cmd/server

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

run:
	$(GORUN) $(MAIN_PACKAGE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

docker:
	docker build -t banking-service:latest .

# Database migrations
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

# Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $${name}

# Start development database
db-start:
	docker run --name banking-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=banking -p 5432:5432 -d postgres:14

# Stop development database
db-stop:
	docker stop banking-postgres
	docker rm banking-postgres