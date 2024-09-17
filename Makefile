.PHONY: lint
lint: bootstrap
	golangci-lint run --max-same-issues 0 --timeout 10m

.PHONY: swagger
swagger: SHELL:=/bin/bash
swagger: bootstrap
	swagger generate spec -o ./swaggerui/html/swagger.json --scan-models --exclude-deps
	swagger validate ./swaggerui/html/swagger.json

HAS_SWAGGER       := $(shell command -v swagger;)
HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)

bootstrap:
ifndef HAS_SWAGGER
	go install github.com/go-swagger/go-swagger/cmd/swagger@v0.31.0
endif
ifndef HAS_GOLANGCI_LINT
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3
endif