APP_NAME    := envdoc
MODULE      := github.com/tendant/envdoc
BIN_DIR     := bin
IMAGE_NAME  := $(APP_NAME)
IMAGE_TAG   := latest
COVERAGE    := coverage.out

.DEFAULT_GOAL := help

.PHONY: build run test test-cover lint vet clean docker-build deploy help

## build: compile the CLI binary
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/envdoc

## run: run locally (use RULES=path/to/rules.yaml to add validation rules)
run: build
ifdef RULES
	$(BIN_DIR)/$(APP_NAME) -rules $(RULES)
else
	$(BIN_DIR)/$(APP_NAME)
endif

## test: run all tests
test:
	go test ./... -v

## test-cover: run tests with coverage report
test-cover:
	go test ./... -coverprofile=$(COVERAGE)
	go tool cover -func=$(COVERAGE)

## vet: run go vet
vet:
	go vet ./...

## lint: run vet (add golangci-lint if available)
lint: vet
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint not installed, skipping"

## clean: remove build artifacts
clean:
	rm -rf $(BIN_DIR) $(COVERAGE)

## docker-build: build Docker image
docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

## deploy: build and push Docker image
deploy: docker-build
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

## help: show this help
help:
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/^## //'
