PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)
ifndef ARCH
ARCH=$(shell uname -m)
endif
ifeq ($(ARCH),x86_64)
	ARCH=amd64
endif
ifeq ($(ARCH),aarch64)
	ARCH=arm64
endif

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

# Buildx on GitHub-hosted runners attaches provenance attestations by default.
# That makes ${VERSION}.amd64 a manifest list, and `docker manifest create`
# fails with "is a manifest list". Keep each arch tag as a single image, and
# combine them with imagetools so ${VERSION} is linux/amd64 + linux/arm64.
DOCKER_BUILD_FLAGS ?= --provenance=false --sbom=false

docker: clean docker-hub docker-openshift docker-fedtest
	docker tag moov/fed:$(VERSION).$(ARCH) moov/fed:$(VERSION)
	docker tag moov/fed:$(VERSION).$(ARCH) moov/fed:latest
	docker tag quay.io/moov/fed:$(VERSION).$(ARCH) quay.io/moov/fed:$(VERSION)
	docker tag quay.io/moov/fed:$(VERSION).$(ARCH) quay.io/moov/fed:latest
	docker tag moov/fedtest:$(VERSION).$(ARCH) moov/fedtest:$(VERSION)
	docker tag moov/fedtest:$(VERSION).$(ARCH) moov/fedtest:latest

.PHONY: docker-hub docker-hub-push docker-hub-manifest
docker-hub:
	docker build --pull ${DOCKER_BUILD_FLAGS} -t moov/fed:$(VERSION).$(ARCH) -f Dockerfile .

docker-hub-push:
	docker push moov/fed:$(VERSION).$(ARCH)

docker-hub-manifest:
	@if docker manifest inspect moov/fed:$(VERSION) >/dev/null 2>&1; then \
		echo "error: docker tag already exists: moov/fed:$(VERSION)" >&2; \
		exit 1; \
	fi
	docker buildx imagetools create -t moov/fed:$(VERSION) moov/fed:$(VERSION).arm64 moov/fed:$(VERSION).amd64
	docker buildx imagetools create -t moov/fed:latest moov/fed:$(VERSION).arm64 moov/fed:$(VERSION).amd64

.PHONY: docker-openshift docker-openshift-push docker-openshift-manifest
docker-openshift:
	docker build --pull ${DOCKER_BUILD_FLAGS} --build-arg VERSION=$(VERSION) -t quay.io/moov/fed:$(VERSION).$(ARCH) -f Dockerfile-openshift .

docker-openshift-push:
	docker push quay.io/moov/fed:$(VERSION).$(ARCH)

docker-openshift-manifest:
	@if docker manifest inspect quay.io/moov/fed:$(VERSION) >/dev/null 2>&1; then \
		echo "error: docker tag already exists: quay.io/moov/fed:$(VERSION)" >&2; \
		exit 1; \
	fi
	docker buildx imagetools create -t quay.io/moov/fed:$(VERSION) quay.io/moov/fed:$(VERSION).arm64 quay.io/moov/fed:$(VERSION).amd64
	docker buildx imagetools create -t quay.io/moov/fed:latest quay.io/moov/fed:$(VERSION).arm64 quay.io/moov/fed:$(VERSION).amd64

.PHONY: docker-fedtest docker-fedtest-push docker-fedtest-manifest
docker-fedtest:
	docker build --pull ${DOCKER_BUILD_FLAGS} -t moov/fedtest:$(VERSION).$(ARCH) -f Dockerfile-fedtest .

docker-fedtest-push:
	docker push moov/fedtest:$(VERSION).$(ARCH)

docker-fedtest-manifest:
	@if docker manifest inspect moov/fedtest:$(VERSION) >/dev/null 2>&1; then \
		echo "error: docker tag already exists: moov/fedtest:$(VERSION)" >&2; \
		exit 1; \
	fi
	docker buildx imagetools create -t moov/fedtest:$(VERSION) moov/fedtest:$(VERSION).arm64 moov/fedtest:$(VERSION).amd64
	docker buildx imagetools create -t moov/fedtest:latest moov/fedtest:$(VERSION).arm64 moov/fedtest:$(VERSION).amd64

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
