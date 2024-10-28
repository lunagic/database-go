.PHONY: full clean lint lint-npm lint-go fix fix-npm fix-go test test-npm test-go build watch docs-go

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
lint: lint-npm lint-go

lint-npm:
	npm install
	npm run lint

lint-go:
	go get -d ./...
	go mod tidy
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55

## Fix the project
fix: fix-npm fix-go

fix-npm:
	npm install
	npm run fix

fix-go:
	go mod tidy
	gofmt -s -w .

## Test the project
test: test-npm test-go

test-npm:
	npm install
	npm run test

test-go:
	@mkdir -p .stencil/tmp/coverage/go/
	@go install github.com/boumenot/gocover-cobertura@latest
	go test -p 1 -count=1 -cover -coverprofile .stencil/tmp/coverage/go/profile.txt ./...
	@go tool cover -func .stencil/tmp/coverage/go/profile.txt | awk '/^total/{print $$1 " " $$3}'
	@go tool cover -html .stencil/tmp/coverage/go/profile.txt -o .stencil/tmp/coverage/go/coverage.html
	@gocover-cobertura < .stencil/tmp/coverage/go/profile.txt > .stencil/tmp/coverage/go/cobertura-coverage.xml

## Build the project
build:

## Watch the project
watch:

## Run the docs server for the project
docs-go:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo "listening on http://127.0.0.1:6060/pkg/github.com/lunagic/database-go"
	@godoc -http=127.0.0.1:6060
