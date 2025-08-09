.PHONY: test build run migrate-up migrate-down docker-build docker-up docker-down

# Test configuration
TEST_DATABASE_DSN ?= postgres://postgres:postgres@localhost:5432/egonez_test?sslmode=disable

# Application configuration
BOT_TOKEN ?= your_bot_token_here

# Test the application
test:
	TEST_DATABASE_DSN=$(TEST_DATABASE_DSN) go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Build the application
build:
	go build -o tmp/main ./cmd/bot

# Run the application
run:
	BOT_TOKEN=$(BOT_TOKEN) \
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=postgres \
	DB_PASSWORD=postgres \
	DB_NAME=egonez \
	DB_SSLMODE=disable \
	LOG_LEVEL=debug \
	go run ./cmd/bot

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	@docker-compose exec -T postgres psql -U postgres -d egonez -c "$(MIGRATION_SQL)"

# Rollback database migrations
migrate-down:
	@echo "Rolling back database migrations..."
	@docker-compose exec -T postgres psql -U postgres -d egonez -c "$(ROLLBACK_SQL)"

# Build Docker image
docker-build:
	docker-compose build

# Start the application with Docker
docker-up:
	docker-compose up -d

# Stop the application
docker-down:
	docker-compose down

# Show logs
docker-logs:
	docker-compose logs -f
