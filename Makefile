.PHONY: build clean development and deploy
run:
	@echo "Running application"
	go run cmd/http/api.go

migrate-up:
	@echo "Running migration"
	migrate -database "mysql://root:123456@tcp(localhost:3306)/transaction" -path db/migrations up

migrate-down:
	@echo "Rollback migration"
	migrate -database "mysql://root:123456@tcp(localhost:3306)/transaction" -path db/migrations down

install:
	@echo "Running download lib"
	go mod download

dev:
	@echo "Running the application"
	go run -tags dynamic cmd/http/api.go

unit-test:
	@echo "Running tests"
	mkdir -p ./test/coverage && \
		CGO_ENABLED=1 GOOS=linux go test $(BUILD_ARGS) -v ./... -coverprofile=./test/coverage/coverage.out

test-dev:
	@echo "Running tests"
	mkdir -p ./test/coverage && \
		CGO_ENABLED=1 go test -tags dynamic -v ./... -coverprofile=./test/coverage/coverage.out

coverage:
	@echo "Running tests with coverage"
	go tool cover -html=./test/coverage/coverage.out

build:
	@echo "Building the application"
	CGO_ENABLED=1 GOOS=linux go build $(BUILD_ARGS) -a -o build/bin/main cmd/http/api.go

start:
	@echo "Start the application"
	./build/bin/main

lint:
	@echo "Running linter"
	golangci-lint run

scan:
	@echo "Running security scan"
	gosec ./...