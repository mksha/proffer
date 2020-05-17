include .env

default: fmt fmt-check mode-check tidy tidy-check

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

build:
	@go build -o bin/proffer .

validate:
	bin/proffer validate test/proffer.yml

apply:
	bin/proffer apply test/proffer.yml

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go

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

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
