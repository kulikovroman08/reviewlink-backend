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

migrate-up-test:
	@set -a; source .env.test; set +a; \
	migrate -path ./migrations -database "$$DB_URL_TEST" up

migrate-down-test:
	@set -a; source .env.test; set +a; \
	migrate -path ./migrations -database "$$DB_URL_TEST" down

migrate-drop-test:
	@set -a; source .env.test; set +a; \
	migrate -path ./migrations -database "$$DB_URL_TEST" drop -f

test-integration:
	@PROJECT_ROOT=$$(pwd) APP_ENV=test go test ./internal/tests/integration/... -v -p 1 -count 1

