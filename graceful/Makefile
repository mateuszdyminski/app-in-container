GOLANG_VERSION := 1.13.8
ALPINE_VERSION := 3.11

VERSION ?= $(shell git describe --always)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
LAST_COMMIT_USER ?= $(shell git log -1 --format='%cn <%ce>')
LAST_COMMIT_HASH ?= $(shell git log -1 --format=%H)
LAST_COMMIT_TIME ?= $(shell git log -1 --format=%cd --date=format:'%Y-%m-%d %H:%M:%S')

NAME := app-container-graceful
DOCKERHUB_REPO := mateuszdyminski
GIT_REPO := github.com/mateuszdyminski/app-in-container/graceful

.DEFAULT_GOAL := build
.PHONY: build docker-build docker-push release run help

build: ## Builds App locally
	CGO_ENABLED=0 \
	go build \
	-v \
	-ldflags "-s -w \
	-X 'main.AppVersion=$(VERSION)' \
	-X 'main.GitVersion=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.LastCommitUser=$(LAST_COMMIT_USER)' \
	-X 'main.LastCommitHash=$(LAST_COMMIT_HASH)' \
	-X 'main.LastCommitTime=$(LAST_COMMIT_TIME)'" \
	-mod vendor -v -o $(NAME)-$(VERSION) .

docker-build: ## Builds Docker image with App
	docker build \
	--pull \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg NAME="$(NAME)" \
	--build-arg VERSION="$(VERSION)" \
	--build-arg BUILD_TIME="$(BUILD_TIME)" \
	--build-arg LAST_COMMIT_USER="$(LAST_COMMIT_USER)" \
	--build-arg LAST_COMMIT_HASH="$(LAST_COMMIT_HASH)" \
	--build-arg LAST_COMMIT_TIME="$(LAST_COMMIT_TIME)" \
	--label="build.version=$(VERSION)" \
	--tag="$(DOCKERHUB_REPO)/$(NAME):$(VERSION)" \
	--tag="$(DOCKERHUB_REPO)/$(NAME):latest" \
	.

docker-push: ## Pushes current build version of Docker image to the registry
	docker push "$(DOCKERHUB_REPO)/$(NAME):$(VERSION)"
	docker push "$(DOCKERHUB_REPO)/$(NAME):latest"

release: docker-build docker-push ## Combines `docker-build` and `docker-push` commands

run: ## Runs App in development mode locally
	APP_HTTP_PORT="8080" \
	APP_DB_HOST="localhost" \
	APP_DB_PORT=3306 \
	APP_DB_USER="root" \
	APP_DB_PASSWORD="password" \
	APP_DB_DBNAME="users" \
	go run .

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
