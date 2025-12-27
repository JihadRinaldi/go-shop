.PHONY: build run dev lint docker-up docker-down

build:
	go build -o bin/app ./cmd/api

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint:
	golangci-lint run

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down