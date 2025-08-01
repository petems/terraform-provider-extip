TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TARGETS=darwin linux windows
TESTARGS?=-race -coverprofile=profile.out -covermode=atomic
LINT := golangci-lint

default: build

build: fmtcheck
	go install

targets: $(TARGETS)

$(TARGETS):
	GOOS=$@ GOARCH=amd64 go build -o "dist/terraform-provider-extip${TRAVIS_TAG}_$@_amd64"
	zip -j dist/terraform-provider-extip${TRAVIS_TAG}_$@_amd64.zip dist/terraform-provider-extip${TRAVIS_TAG}_$@_amd64

## Run linter
lint:
	@echo "Checking golangci-lint version..."
	@$(LINT) version | grep -q "golangci-lint has version" || (echo "golangci-lint not found. Please install it first." && exit 1)
	@$(LINT) version | grep -oE "version [0-9]+\.[0-9]+\.[0-9]+" | cut -d' ' -f2 | awk -F. '{if ($$1 > 2 || ($$1 == 2 && $$2 >= 3)) exit 0; else exit 1}' || (echo "golangci-lint version 2.3.0 or higher required. Current version:" && $(LINT) version && exit 1)
	$(LINT) run --config=.golangci.yml

lint-fix:
	$(LINT) run --fix --config=.golangci.yml

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

# Fast tests (skip real network tests)
test-short: fmtcheck
	go test -short -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -short $(TESTARGS) -timeout=30s -parallel=4

# Verbose tests
test-verbose: fmtcheck
	go test -v $(TEST)

# Fast verbose tests
test-verbose-short: fmtcheck
	go test -short -v $(TEST)

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: cover_html
cover_html: test ## Runs go test with coverage
	@go tool cover -html=profile.out

# DevOps targets
.PHONY: ci-setup
ci-setup: ## Setup CI environment
	@echo "Setting up CI environment..."
	@go mod download
	@go mod verify

.PHONY: security-scan
security-scan: ## Run security scans
	@echo "Running security scans..."
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "Installing gosec..."; \
		go install github.com/securecode/gosec/v2/cmd/gosec@latest; \
	fi
	@gosec ./...
	@if ! command -v govulncheck >/dev/null 2>&1; then \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	@govulncheck ./...

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@go mod vendor

.PHONY: pre-commit
pre-commit: fmtcheck vet lint test-short ## Run pre-commit checks
	@echo "Pre-commit checks passed!"

.PHONY: ci-test
ci-test: ci-setup pre-commit test ## Full CI test suite
	@echo "CI tests completed!"

.PHONY: release-build
release-build: ## Build release binaries
	@echo "Building release binaries..."
	@goreleaser build --snapshot --rm-dist

.PHONY: release-dry-run
release-dry-run: ## Dry run release process
	@echo "Running release dry run..."
	@goreleaser release --snapshot --rm-dist --skip-publish

.PHONY: docs-generate
docs-generate: ## Generate provider documentation
	@echo "Generating documentation..."
	@if ! command -v tfplugindocs >/dev/null 2>&1; then \
		echo "Installing tfplugindocs..."; \
		go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest; \
	fi
	@tfplugindocs generate

.PHONY: docs-validate
docs-validate: ## Validate provider documentation
	@echo "Validating documentation..."
	@if ! command -v tfplugindocs >/dev/null 2>&1; then \
		echo "Installing tfplugindocs..."; \
		go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest; \
	fi
	@tfplugindocs validate

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build test test-short test-verbose test-verbose-short testacc vet fmt fmtcheck errcheck vendor-status test-compile
.PHONY: ci-setup security-scan deps-update pre-commit ci-test release-build release-dry-run docs-generate docs-validate help
