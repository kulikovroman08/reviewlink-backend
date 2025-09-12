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

test-integration:
	@APP_ENV=test DB_URL="postgres://app:secret@localhost:5432/reviewlink_test?sslmode=disable" go test ./internal/tests/... -v -p 1 -count 1


