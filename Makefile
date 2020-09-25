all: test check-coverage lint build-examples ## test, check coverage, lint and build examples
	@if [ -e .git/rebase-merge ]; then git --no-pager log -1 --pretty='%h %s'; fi
	@echo '$(COLOUR_GREEN)Success$(COLOUR_NORMAL)'

clean::  ## Remove generated files

.PHONY: all clean

# -- Test --------------------------------------------------------------

COVERFILE = coverage.out
COVERAGE = 82.0

test:  ## Run tests and generate a coverage file
	go test -covermode=atomic -coverprofile=$(COVERFILE) -race ./...
	go mod tidy

integration: ## Run integration tests which interact over the network
	go test -covermode=atomic -tags integration -race ./...

bench: ## Run benchmarks only
	go test -run=NONE -bench=. ./...

check-coverage: test  ## Check that test coverage meets the required level
	@go tool cover -func=$(COVERFILE) | $(CHECK_COVERAGE) || $(FAIL_COVERAGE)

cover: test  ## Show test coverage in your browser
	go tool cover -html=$(COVERFILE)

clean::
	rm -f $(COVERFILE)

CHECK_COVERAGE = awk -F '[ \t%]+' '/^total:/ {print; if ($$3 < $(COVERAGE)) exit 1}'
FAIL_COVERAGE = { echo '$(COLOUR_RED)FAIL - Coverage below $(COVERAGE)%$(COLOUR_NORMAL)'; exit 1; }

.PHONY: test check-coverage cover

# -- Lint --------------------------------------------------------------

GOLINT_VERSION = 1.24.0
GOLINT_INSTALLED_VERSION = $(or $(word 4,$(shell golangci-lint --version 2>/dev/null)),0.0.0)
GOLINT_MIN_VERSION = $(shell printf '%s\n' $(GOLINT_VERSION) $(GOLINT_INSTALLED_VERSION) | sort -V | head -n 1)

lint: ## Lint source code
ifeq ($(GOLINT_MIN_VERSION), $(GOLINT_VERSION))
	golangci-lint run
else
lint: lint-with-docker
endif

lint-with-docker:
	docker run --rm -v $(PWD):/src -w /src golangci/golangci-lint:v$(GOLINT_VERSION) golangci-lint run

.PHONY: lint

# -- Build examples ----------------------------------------------------

# The `_examples` directory by starting with `_` is ignored by the Go
# tools. We want to ignore it so we don't create empty godocs for it,
# but we also want to ensure that it still builds so add the special
# target `build-examples`:
build-examples:
	go build -o log-examples ./log/_examples

.PHONY: build-examples

# -- Generate ----------------------------------------------------------

generate:
	go generate ./...

.PHONY: generate

# -- Utilities ---------------------------------------------------------

COLOUR_NORMAL = $(shell tput sgr0 2>/dev/null)
COLOUR_RED    = $(shell tput setaf 1 2>/dev/null)
COLOUR_GREEN  = $(shell tput setaf 2 2>/dev/null)
COLOUR_WHITE  = $(shell tput setaf 7 2>/dev/null)

help:
	@awk -F ':.*## ' 'NF == 2 && $$1 ~ /^[A-Za-z0-9_-]+$$/ { printf "$(COLOUR_WHITE)%-30s$(COLOUR_NORMAL)%s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: help
