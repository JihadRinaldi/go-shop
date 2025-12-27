DB_URL=postgresql://postgres:password@localhost:5432/go_shop?sslmode=disable

.PHONY: build run dev lint docker-up docker-down resetdb new_migration migrateup migratedown migrateup_n migratedown_n

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

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

resetdb: migratedown migrateup

migrateup_n:
	@read -p "Enter number of migrations: " n; \
	migrate -path db/migrations -database "$(DB_URL)" -verbose up $$n

migratedown_n:
	@read -p "Enter number of migrations: " n; \
	migrate -path db/migrations -database "$(DB_URL)" -verbose down $$n