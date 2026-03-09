.PHONY: build run test update clean help

SHELL := /bin/bash

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build and deploy (full pipeline)
	@bash ../run.sh

run: ## Build, deploy, and run with optional ARGS (make run ARGS="list --json")
	@bash ../run.sh -r $(ARGS)

test: ## Run tests with report
	@bash ../run.sh -t

update: ## Pull, build, deploy, and sync PATH binary
	@bash ../run.sh -u

no-pull: ## Build without git pull
	@bash ../run.sh -n

no-deploy: ## Build without deploying
	@bash ../run.sh -d

clean: ## Remove build artifacts
	@rm -rf bin/
	@echo "Build artifacts removed."
