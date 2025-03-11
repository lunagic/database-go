.PHONY: full clean lint lint-go fix fix-go test test-go build watch docs-go

SHELL=/bin/bash -o pipefail
$(shell git config core.hooksPath ops/git-hooks)
PROJECT_NAME := $(shell basename $(CURDIR))
GO_PATH := $(shell go env GOPATH 2> /dev/null)
PATH := /usr/local/bin:$(GO_PATH)/bin:$(PATH)

full: clean lint test build

## Clean the project of temporary files
clean:
	git clean -Xdff --exclude="!.env*local"

## Lint the project
lint: lint-go

lint-go:
	go get ./...
	go mod tidy
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55

## Fix the project
fix: fix-go

fix-go:
	go mod tidy
	gofmt -s -w .

## Test the project
test: test-go

test-go:
	@mkdir -p tmp/coverage/go/
	@go install github.com/boumenot/gocover-cobertura@latest
	go test -cover -coverprofile tmp/coverage/go/profile.txt ./...
	@go tool cover -func tmp/coverage/go/profile.txt | awk '/^total/{print $$1 " " $$3}'
	@go tool cover -html tmp/coverage/go/profile.txt -o tmp/coverage/go/coverage.html
	@gocover-cobertura < tmp/coverage/go/profile.txt > tmp/coverage/go/cobertura-coverage.xml

## Build the project
build:

## Watch the project
watch:

## Run the docs server for the project
docs-go:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo "listening on http://127.0.0.1:6060/pkg/github.com/lunagic/database-go"
	@godoc -http=127.0.0.1:6060
