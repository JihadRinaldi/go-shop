DB_URL=postgresql://postgres:password@localhost:5432/go_shop?sslmode=disable

.PHONY: build run dev lint docker-up docker-down resetdb new_migration migrateup migratedown migrateup_n migratedown_n test mocks clean-mocks

build:
	@echo "Building all binaries...."
	@mkdir -p bin
	@for cmd in cmd/*/; do \
    		if [ -d "$$cmd" ]; then \
    			binary=$$(basename $$cmd); \
    			echo "Building $$binary..."; \
    			go build -o bin/$$binary ./$$cmd; \
    		fi \
    	done

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

# Test targets
test:
	@echo "Running tests..."
	go test -v -cover ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Mock generation targets
mocks:
	@echo "Generating mocks..."
	@mockery

clean-mocks:
	@echo "Cleaning mocks..."
	@rm -rf internal/mocks
	@echo "Mocks cleaned"

regenerate-mocks: clean-mocks mocks
	@echo "Mocks regenerated successfully"