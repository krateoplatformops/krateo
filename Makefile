# Set the shell to bash always
SHELL := /bin/bash

# Look for a .env file, and if present, set make variables from it.
ifneq (,$(wildcard ../../.env))
	include ../../.env
	export $(shell sed 's/=.*//' ../../.env)
endif

NAME := krateoctl
ORG := krateoplatformops
ORG_REPO := $(ORG)/$(NAME)
ROOT_PACKAGE := github.com/$(ORG_REPO)
# set dev version unless VERSION is explicitly set via environment
VERSION ?= $(shell echo "$$(git for-each-ref refs/tags/ --count=1 --sort=-version:refname --format='%(refname:short)' 2>/dev/null)-dev+$(REV)" | sed 's/^v//')

KIND_CLUSTER_NAME ?= local-dev
KUBECONFIG ?= $(HOME)/.kube/config

# Tools
KIND=$(shell which kind)
KUBECTL=$(shell which kubectl)
HELM=$(shell which helm)

.DEFAULT_GOAL := dev

.PHONY: dev
dev: ## dev build
dev: clean tools generate vet fmt lint test mod-tidy

.PHONY: ci
ci: ## CI build
ci: dev

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
	rm -rf dist
	rm -f coverage.*

.PHONY: tools
tools: ## go install tools
	$(call print-target)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser@latest

.PHONY: generate
generate: ## go generate
	$(call print-target)
	go generate ./...

.PHONY: vet
vet: ## go vet
	$(call print-target)
	go vet ./...

.PHONY: fmt
fmt: ## go fmt
	$(call print-target)
	go fmt ./...

.PHONY: lint
lint: ## golangci-lint
	$(call print-target)
	golangci-lint run

.PHONY: test
test: ## go test with race detector and code covarage
	$(call print-target)
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: mod-tidy
mod-tidy: ## go mod tidy
	$(call print-target)
	go mod tidy

.PHONY: diff
diff: ## git diff
	$(call print-target)
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi

.PHONY: build
build: ## goreleaser --snapshot --skip-publish --rm-dist
build: tools
	$(call print-target)
	ROOT_PACKAGE=github.com/$(ORG_REPO) goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release: ## goreleaser --rm-dist
release: tools
	$(call print-target)
	ROOT_PACKAGE=github.com/$(ORG_REPO) goreleaser --rm-dist

.PHONY: run
run: ## go run
	@go run -race .

.PHONY: go-clean
go-clean: ## go clean build, test and modules caches
	$(call print-target)
	go clean -r -i -cache -testcache -modcache

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: kind.up
kind.up: ## starts a KinD cluster for local development
	@$(KIND) get kubeconfig --name $(KIND_CLUSTER_NAME) >/dev/null 2>&1 || $(KIND) create cluster --name=$(KIND_CLUSTER_NAME)

.PHONY: kind.down
kind.down: ## shuts down the KinD cluster
	@$(KIND) delete cluster --name=$(KIND_CLUSTER_NAME)

.PHONY: kind.clean
kind.clean:
	@$(HELM) uninstall argocd --namespace argo-system || true
	@$(HELM) uninstall crossplane --namespace crossplane-system || true

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
