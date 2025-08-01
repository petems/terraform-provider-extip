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

.PHONY: build test test-short test-verbose test-verbose-short testacc vet fmt fmtcheck errcheck vendor-status test-compile
