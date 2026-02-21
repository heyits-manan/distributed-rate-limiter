.PHONY: build run test lint clean docker-build docker-up docker-down

APP_NAME := rate-limiter
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run: build
	./$(BUILD_DIR)/$(APP_NAME)

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)

docker-build:
	docker build -t $(APP_NAME) .

docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down
