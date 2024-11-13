APP_NAME ?= supotsu-backend
CURRENT_DIR := $(shell basename $(shell pwd))

TARGET_OS ?= $(shell go env GOOS)
TARGET_ARCH ?= $(shell go env GOARCH)

.PHONY: docker-local
docker-local: build
	@docker build . -t $(APP_NAME):local

.PHONY: build
build:
	@CGO_ENABLED=0 go build -o dist/backend-$(TARGET_OS)-$(TARGET_ARCH) ./cmd/app

.PHONY: up
up:
	@docker compose up -d

.PHONY: down
down:
	@docker compose down

.PHONY: clean
clean:
	@docker volume rm $(CURRENT_DIR)_log $(CURRENT_DIR)_config $(CURRENT_DIR)_pb-data
	@rm -rf dist
