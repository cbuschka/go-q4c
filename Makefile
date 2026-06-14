PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
TOOLS_DIR := $(PROJECT_DIR)/.gotools
SHELL = /bin/bash

LINT_VER := v2.12.2

.PHONY:	all
all:	test

.PHONY:	test
test:
	go test ./...

.PHONY:	update-deps
update-deps:
	go get -d -v -u all

.PHONY:	list-deps
list-deps:
	go list -m -u all

.PHONY: lint-install
lint-install:
	@echo "Checking golangci-lint..." && \
	[[ -f ./.bin/golangci-lint ]] || (curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./.bin $(LINT_VER))

.PHONY: lint
lint: lint-install
	mkdir -p $(TOOL_DIR)/bin && \
	./.bin/golangci-lint run ./...