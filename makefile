PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

.PHONY: build docker release check

build:
# main FED binary
	CGO_ENABLED=0 go build -o ./bin/server github.com/moov-io/fed/cmd/server
# fedtest binary
	CGO_ENABLED=0 go build -o bin/fedtest ./cmd/fedtest

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=40.0 ./lint-project.sh
endif

.PHONY: client
client:
ifeq ($(OS),Windows_NT)
	@echo "Please generate client on macOS or Linux, currently unsupported on windows."
else
# Versions from https://github.com/OpenAPITools/openapi-generator/releases
	@chmod +x ./openapi-generator
	@rm -rf ./client
	OPENAPI_GENERATOR_VERSION=4.3.1 ./openapi-generator generate --git-user-id=moov-io --git-repo-id=fed --package-name client -i ./api/client.yaml -g go -o ./client
	rm -f ./client/go.mod ./client/go.sum ./client/.travis.yml ./client/git_push.sh
	go fmt ./...
	go build github.com/moov-io/fed/client
	go test ./client/...
endif

.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	@echo "Skipping cleanup on Windows, currently unsupported."
else
	@rm -rf ./bin/ cover.out coverage.txt openapi-generator-cli-*.jar misspell* staticcheck* lint-project.sh
endif

dist: clean client build
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=1 GOOS=windows go build -o bin/fed.exe github.com/moov-io/fed/cmd/server
else
	CGO_ENABLED=1 GOOS=$(PLATFORM) go build -o bin/fed-$(PLATFORM)-amd64 github.com/moov-io/fed/cmd/server
endif

docker: clean
# main FED image
	docker build --pull -t moov/fed:$(VERSION) -f Dockerfile .
	docker tag moov/fed:$(VERSION) moov/fed:latest
# OpenShift Docker image
	docker build --pull -t quay.io/moov/fed:$(VERSION) -f Dockerfile-openshift --build-arg VERSION=$(VERSION) .
	docker tag quay.io/moov/fed:$(VERSION) quay.io/moov/fed:latest
# fedtest image
	docker build --pull -t moov/fedtest:$(VERSION) -f Dockerfile-fedtest ./
	docker tag moov/fedtest:$(VERSION) moov/fedtest:latest

clean-integration:
	docker compose kill
	docker compose rm -v -f

test-integration: clean-integration
	docker compose up -d
	sleep 5
	./bin/fedtest -local

release: docker AUTHORS
	go vet ./...
	go test -coverprofile=cover-$(VERSION).out ./...
	git tag -f $(VERSION)

release-push:
	docker push moov/fed:$(VERSION)
	docker push moov/fed:latest

quay-push:
	docker push quay.io/moov/fed:$(VERSION)
	docker push quay.io/moov/fed:latest

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
