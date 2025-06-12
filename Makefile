PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

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
