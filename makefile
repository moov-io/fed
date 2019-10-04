PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

.PHONY: build docker release check

build: check
# main FED binary
	CGO_ENABLED=1 go build -o ./bin/server github.com/moov-io/fed/cmd/server
# fedtest binary
	CGO_ENABLED=0 go build -o bin/fedtest ./cmd/fedtest

check:
	go fmt ./...
	@mkdir -p ./bin/

.PHONY: client
client:
# Versions from https://github.com/OpenAPITools/openapi-generator/releases
	@chmod +x ./openapi-generator
	@rm -rf ./client
	OPENAPI_GENERATOR_VERSION=4.1.3 ./openapi-generator generate -i openapi.yaml -g go -o ./client
	rm -f client/go.mod client/go.sum
	go fmt ./...
	go build github.com/moov-io/fed/client
	go test ./client

.PHONY: clean
clean:
	@rm -rf client/
	@rm -rf bin/
	@rm -f openapi-generator-cli-*.jar

dist: clean client build
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=1 GOOS=windows go build -o bin/fed-windows-amd64.exe github.com/moov-io/fed/cmd/server
else
	CGO_ENABLED=1 GOOS=$(PLATFORM) go build -o bin/fed-$(PLATFORM)-amd64 github.com/moov-io/fed/cmd/server
endif

docker:
# main FED image
	docker build --pull -t moov/fed:$(VERSION) -f Dockerfile .
	docker tag moov/fed:$(VERSION) moov/fed:latest
# fedtest image
	docker build --pull -t moov/fedtest:$(VERSION) -f Dockerfile-fedtest ./
	docker tag moov/fedtest:$(VERSION) moov/fedtest:latest

release: docker AUTHORS
	go vet ./...
	go test -coverprofile=cover-$(VERSION).out ./...
	git tag -f $(VERSION)

release-push:
	docker push moov/fed:$(VERSION)
	docker push moov/fed:latest

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
