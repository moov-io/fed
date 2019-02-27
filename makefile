PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

.PHONY: build build-server build-examples docker release check

build-server:
	CGO_ENABLED=1 go build -o ./bin/server github.com/moov-io/fed/cmd/server

check:
	go fmt ./...
	@mkdir -p ./bin/

.PHONY: client
client:
# Versions from https://github.com/OpenAPITools/openapi-generator/releases
	@chmod +x ./openapi-generator
	@rm -rf ./client
	OPENAPI_GENERATOR_VERSION=4.0.0-beta2 ./openapi-generator generate -i openapi.yaml -g go -o ./client
	go fmt ./client
	go build github.com/moov-io/fed/client
	go test ./client

.PHONY: clean
clean:
	@rm -rf client/
	@rm -f openapi-generator-cli-*.jar

dist: clean client build
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=1 GOOS=windows go build -o bin/fed-windows-amd64.exe github.com/moov-io/fed/cmd/server
else
	CGO_ENABLED=1 GOOS=$(PLATFORM) go build -o bin/fed-$(PLATFORM)-amd64 github.com/moov-io/fed/cmd/server
endif

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...
cover-web:
	go tool cover -html=cover.out

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@
