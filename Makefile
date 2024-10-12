# Variables
DB_USER=postgres
DB_NAME=erp
MIGRATION_FILE=db/migration.sql
APP_BINARY=main.go
TEST_PATTERN=./...

# Declare targets as PHONY
.PHONY: run migrate test test-race test-cover clean help

# Default target to run the application
run:
	@go run $(APP_BINARY)

# Target to run migrations
migrate:
	@psql -U $(DB_USER) -d $(DB_NAME) -f $(MIGRATION_FILE)

# Target to run tests
test:
	@go test $(TEST_PATTERN) -v

# Target to run tests with race condition detection
test-race:
	@go test $(TEST_PATTERN) -race -v

# Target to run tests with code coverage report
test-cover:
	@go test $(TEST_PATTERN) -cover -v

# Target to clean any generated files (if needed)
clean:
	@echo "Nothing to clean for now..."

# Target to show help
help:
	@echo "Makefile commands:"
	@echo "  run         - Run the application"
	@echo "  migrate     - Run database migration"
	@echo "  test        - Run tests"
	@echo "  test-race   - Run tests with race condition detection"
	@echo "  test-cover  - Run tests with coverage report"
	@echo "  clean       - Clean up generated files"
	@echo "  help        - Show this help message"
