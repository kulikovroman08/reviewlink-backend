SHELL := /bin/bash

.PHONY: lint

lint:
	golangci-lint run ./...

migrate-up:
	@set -a; source .env; set +a; \
	migrate -path ./migrations -database "$$DB_URL" up

migrate-down:
	@set -a; source .env; set +a; \
	migrate -path ./migrations -database "$$DB_URL" down

create-migration:
	migrate create -ext sql -dir ./migrations $(name)

migrate-force-zero:
	@set -a; source .env; set +a; \
	migrate -path ./migrations -database "$$DB_URL" force 0

migrate-force-zero-test:
	@set -a; source .env; set +a; \
	migrate -path ./migrations -database "$$DB_URL_TEST" force 0

test-integration:
	@PROJECT_ROOT=$$(pwd) APP_ENV=test go test ./internal/tests/integration/... -v -p 1 -count 1

