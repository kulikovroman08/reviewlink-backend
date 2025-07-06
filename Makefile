SHELL := /bin/bash

.PHONY: lint

lint:
	golangci-lint run ./...