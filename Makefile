GOLANGCI_LINT_VERSION ?= v1.64.6
OSNAME ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
ARTIFACTS_DIR := ./artifacts
DOCKER_TAG := "http-v1"

.PHONY: build
build:
	CGO_ENABLED=0 \
	GOOS=linux GOARCH=amd64 \
	go build \
		-ldflags="-s -w" \
		-o $(ARTIFACTS_DIR)/svc \
		./cmd/main.go

.PHONY: dockerise
dockerise:
	docker build \
	--force-rm \
	--rm \
	-t $(DOCKER_TAG) \
	.

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --allow-parallel-runners

.PHONY: mockgen
mockgen:
	go generate ./...

.PHONY: local-run
local-run:
	$(ARTIFACTS_DIR)/svc

.PHONE: local-compose-start
	docker compose -f docker-composition/postgres.yml up \
		--build \
		--detach

.PHONY: local-compose-stop
	docker compose -f docker-composition/postgres.yml down \
		--remove-orphans

.PHONY: unit_test
unit_test:
	go test -parallel 6 -race -count=1 -coverpkg=./... -coverprofile=unit_coverage.out -v `go list ./... | grep -v /test/`

.PHONY: integration_test
integration_test:
	go test -count=1 -v --tags=integration -coverpkg=./... -coverprofile=int_coverage.out ./test/...
