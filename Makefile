PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY:	all
all:	test

.PHONY:	test
test:
	go test ./...
