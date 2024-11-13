APP_NAME ?= supotsu-backend
CURRENT_DIR := $(shell basename $(shell pwd))

.PHONY: docker-local
docker-local:
	@docker build . -t $(APP_NAME):local

.PHONY: build
build:
	@CGO_ENABLED=0 go build -o $(APP_NAME) ./cmd/app

.PHONY: up
up:
	@docker compose up -d

.PHONY: down
down:
	@docker compose down


.PHONY: clean
clean:
	@docker volume rm $(CURRENT_DIR)_log $(CURRENT_DIR)_config $(CURRENT_DIR)_pb-data
