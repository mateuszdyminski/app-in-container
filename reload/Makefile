GOLANG_VERSION := 1.13.8
ALPINE_VERSION := 3.11

VERSION ?= $(shell git describe --always)
NAME := app-container-reload
DOCKERHUB_REPO := mateuszdyminski

.DEFAULT_GOAL := build
.PHONY: build docker-build docker-push release run help

build: ## Builds App locally
	CGO_ENABLED=0 go build -mod vendor -v -o $(NAME)-$(VERSION) .

docker-build: ## Builds Docker image with App
	docker build \
	--pull \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg NAME="$(NAME)" \
	--build-arg VERSION="$(VERSION)" \
	--label="build.version=$(VERSION)" \
	--tag="$(DOCKERHUB_REPO)/$(NAME):$(VERSION)" \
	--tag="$(DOCKERHUB_REPO)/$(NAME):latest" \
	.

docker-push: ## Pushes current build version of Docker image to the registry
	docker push "$(DOCKERHUB_REPO)/$(NAME):$(VERSION)"
	docker push "$(DOCKERHUB_REPO)/$(NAME):latest"

release: docker-build docker-push ## Combines `docker-build` and `docker-push` commands

run: ## Runs App in development mode locally
	CONFIG_PATH="config" \
	go run .

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
