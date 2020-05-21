include .env

.PHONY: dfault ci fmt fmt-check bsd-mode-check mode-check tidy tidy-check install-lint-deps lint ci-lint vet \
				cover ci-cover test testrace help

default: fmt fmt-check mode-check tidy tidy-check lint

ci: fmt-check mode-check tidy-check ci-lint ## Test in continuous integration

fmt: ## Format Go code
	@go fmt ./...

fmt-check: fmt ## Check go code formatting
	@echo "==> Checking that code complies with go fmt requirements..."
	@git diff --exit-code; if [ $$? -eq 1 ]; then \
		echo "Found files that are not fmt'ed."; \
		echo "You can use the command: \`make fmt\` to reformat code."; \
		exit 1; \
	fi

bsd-mode-check: ## Check that only certain files are executable
	@echo "==> Checking that only certain files are executable..."
	@if [ ! -z "$(BSD_EXECUTABLE_FILES)" ]; then \
		echo "These files should not be executable or they must be white listed in the Makefile:"; \
		echo "$(BSD_EXECUTABLE_FILES)" | xargs -n1; \
		exit 1; \
	else \
		echo "Check passed."; \
	fi

mode-check: ## Check that only certain files are executable
	@echo "==> Checking that only certain files are executable..."
	@if [ ! -z "$(EXECUTABLE_FILES)" ]; then \
		echo "These files should not be executable or they must be white listed in the Makefile:"; \
		echo "$(EXECUTABLE_FILES)" | xargs -n1; \
		exit 1; \
	else \
		echo "Check passed."; \
	fi

install-lint-deps: ## Install linter dependencies
	@echo "==> Updating linter dependencies..."
	@curl -sSfL -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.27.0

lint: install-lint-deps ## Lint Go code
	@if [ ! -z  $(PKG_NAME) ]; then \
		echo "golangci-lint run ./$(PKG_NAME)/..."; \
		golangci-lint run ./$(PKG_NAME)/...; \
	else \
		echo "golangci-lint run ./..."; \
		golangci-lint run ./...; \
	fi

ci-lint: install-lint-deps ## On ci only lint newly added Go source files
	@echo "==> Running linter on newly added Go source files..."
	GO111MODULE=on golangci-lint run --new-from-rev=$(shell git merge-base origin/master HEAD) ./...

vet: ## Vet Go code
	@go vet $(VET)  ; if [ $$? -eq 1 ]; then \
		echo "ERROR: Vet found problems in the code."; \
		exit 1; \
	fi

tidy: ## Remove unused go modules
	@go mod tidy  ; if [ $$? -eq 1 ]; then \
		echo "ERROR: Mod Tidy found problems in the code."; \
		exit 1; \
	fi

tidy-check: tidy ## Check go code for unused modules
	@echo "==> Checking go code for unused modules..."
	@git diff --exit-code; if [ $$? -eq 1 ]; then \
		echo "Found unused modules"; \
		echo "You can use the command: \`make tidy\` to remove unused go modules."; \
		exit 1; \
	fi

# Runs code coverage and open a html page with report
cover:
	go test $(TEST) $(TESTARGS) -timeout=3m -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out
	rm coverage.out

# Runs code coverage and upload the report to codecov.io for ci builds
ci-cover:
	go test $(TEST) $(TESTARGS) -timeout=3m -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out

test: mode-check vet ## Run unit tests
	@go test $(TEST) $(TESTARGS) -timeout=3m

testrace: mode-check vet ## Test with race detection enabled
	@go test $(TEST) $(TESTARGS) -timeout=3m -p=8

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
