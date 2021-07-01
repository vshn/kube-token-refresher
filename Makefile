
IMG_TAG ?= latest
IMG ?= quay.io/vshn/kube-token-refresher:$(IMG_TAG)

all: fmt vet build

.PHONY: build
build: 
	go build

.PHONY: test
test:
	go test ./... -coverprofile cover.out

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: lint
lint: fmt vet ## Invokes the fmt and vet targets
	@echo 'Check for uncommitted changes ...'
	git diff --exit-code

.PHONY: docker-build
docker-build: export GOOS = linux
docker-build: build ## Build the docker image
	docker build .  -t $(IMG) 

.PHONY: help
help: ## Show this help
	@grep -E -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = "(: ).*?## "}; {gsub(/\\:/,":", $$1)}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
