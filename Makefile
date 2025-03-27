SHELL := /bin/bash

.DEFAULT_GOAL := all

# Recommend to `make setup` first

# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./
BINARY_NAME := app
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
	rm -f coverage.*
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -modcache -fuzzcache -x

.PHONY: mod
mod: ## go mod tidy
	$(call print-target)
	go fmt ./...
	go mod tidy -v

.PHONY: setup
setup: ## setup tool, install deps and githook in application
setup:
	chmod +x ./.scripts/*
	./.scripts/setup.sh

.PHONY: upgrade
upgrade: ## upgrade dev tool (golangci-lint,govulncheck,swagger,godoc,pkgsite)
upgrade:
	chmod +x ./.scripts/*
	./.scripts/upgrade-dev-tools.sh

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: all
all: ## mod lint test
all: mod lint test

.PHONY: precommit
precommit: ## validate the branch before commit
precommit: all vuln
	go mod verify
	go vet ./...

.PHONY: ci
ci: ## CI build pipeline
ci: precommit diff

.PHONY: lint
lint: ## golangci-lint
	$(call print-target)
	golangci-lint run --fix

.PHONY: vuln
vuln: ## govulncheck
	$(call print-target)
	govulncheck ./...

.PHONY: test
test: ## go test
	$(call print-target)
	go test -v -race -covermode=atomic -buildvcs -coverpkg=./... ./... | ./.scripts/colorize

.PHONY: coverage
coverage: ## run all tests and display coverage
	go test -v -race -buildvcs -coverprofile=./coverage.out ./...
	go tool cover -html=./coverage.out

.PHONY: diff
diff: ## check if there are any uncommitted changes
	$(call print-target)
	git diff --exit-code
	@RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: docker
docker: ## build the docker image
	@echo "Building the docker image"
	# build and replace the image
	docker build --rm -t api --build-arg GIT_COMMIT=${GIT_COMMIT} .

.PHONY: run
run: ## run the docker container
	@echo "Running the container"
	# run the container
	docker run --env-file ./.env --network host -p 8080:8080 api

.PHONY: swagger
swagger: ## generate swagger file
	swagger generate spec -o ./openapi/swagger.yaml --scan-models

.PHONY: openapi
openapi: ## serve swagger file
	swagger serve -F=swagger --port 8910 ./openapi/swagger.yaml

.PHONY: bump-version
bump-version: ## bump version takes a single argument, the version to bump e.g. v1.0.1
ifeq ($(version),)
	@echo "Error: Missing version argument"
	@echo "Usage: make bump-version version=v0.0.1"
	exit 1
endif
	./.scripts/bump-version.sh ${version}

.PHONY: commit-msg
commit-msg: ## commit-msg: check commit message
	./.scripts/commit-msg.sh

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef

.PHONY: doc
doc: ## doc: show documents
	pkgsite -http 0.0.0.0:6060 -open .
