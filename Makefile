SHELL = bash
.DEFAULT_GOAL := build

.PHONY: build
build:
	CGO_ENABLED=0
	go build \
	-installsuffix cgo \
	-ldflags "-s -w" \
	-o ./bin/sample-api \
	.

.PHONY: test
test:
	go test -cover ./...

.PHONY: lint
lint: bootstrap
	golangci-lint run --max-same-issues 0 --timeout 10m

.PHONY: mocks
mocks: bootstrap
	mockgen -source ./models/controller.go -destination ./models/mock_controller.go -package models

.PHONY: swagger
swagger: SHELL:=/bin/bash
swagger: bootstrap
	swag init -o ./swaggerui/html


HAS_SWAGGER       := $(shell command -v swagger;)
HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)
HAS_MOCKGEN       := $(shell command -v mockgen;)

bootstrap:
ifndef HAS_SWAGGER
	go get -d github.com/swaggo/swag/cmd/swag
endif
ifndef HAS_GOLANGCI_LINT
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3
endif
ifndef HAS_MOCKGEN
	go install github.com/golang/mock/mockgen@v1.6.0
endif