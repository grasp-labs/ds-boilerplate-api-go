.PHONY: lint
lint: bootstrap
	golangci-lint run --max-same-issues 0 --timeout 10m

.PHONY: swagger
swagger: SHELL:=/bin/bash
swagger: bootstrap
	swag init -o ./swaggerui/html

HAS_SWAGGER       := $(shell command -v swagger;)
HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)

bootstrap:
ifndef HAS_SWAGGER
	go get -d github.com/swaggo/swag/cmd/swag
endif
ifndef HAS_GOLANGCI_LINT
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3
endif