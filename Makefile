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