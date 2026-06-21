all: build-site

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "")

LDFLAGS := \
	-X main.Version=$(VERSION) \
	-X main.FullCommitHash=$(COMMIT)

.PHONY: all build run serve tests bench clean install-hooks fmt lint audit build-site

# Points git at the tracked .github/githooks dir so the pre-commit hook
# (fmt, lint, audit, tests, build-check) is active for everyone who builds
# the repo. No-ops quietly outside of a git checkout (e.g. building from a
# release tarball).
install-hooks:
	@git rev-parse --is-inside-work-tree >/dev/null 2>&1 && \
		git config core.hooksPath .github/githooks && \
		echo "anna: pre-commit hook installed (fmt, lint, audit, tests, build)" || true

build: install-hooks
	@echo "anna: building $(VERSION) ($(COMMIT))"
	@go build -ldflags "$(LDFLAGS)"

build-site: build
	@./anna

serve: build
	@./anna -s -p site

tests:
	@echo "anna: running all tests"
	@go test ./...

# Fails and lists offenders if any .go file isn't gofmt'd. Run
# `gofmt -l -w .` yourself to fix them.
fmt:
	@echo "anna: checking gofmt"
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "anna: the following files need gofmt:"; \
		echo "$$unformatted"; \
		echo "anna: run 'gofmt -l -w .' to fix"; \
		exit 1; \
	fi

# go vet always runs; golangci-lint runs too if it's installed, otherwise
# this step is skipped with a hint on how to install it.
# https://golangci-lint.run/welcome/install/
lint:
	@echo "anna: running go vet"
	@go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "anna: running golangci-lint"; \
		golangci-lint run ./...; \
	else \
		echo "anna: golangci-lint not installed, skipping (see https://golangci-lint.run/welcome/install/)"; \
	fi

# Scans module dependencies for known vulnerabilities via govulncheck, if
# installed (requires network access to vuln.go.dev). Skipped otherwise.
# go install golang.org/x/vuln/cmd/govulncheck@latest
audit:
	@echo "anna: running go audit"
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "anna: govulncheck not installed, skipping (see go install golang.org/x/vuln/cmd/govulncheck@latest)"; \
	fi
	
bench:
	@echo "anna: running benchmark"
	@go test -bench . -benchmem -cpuprofile pprof.cpu
	@# to profile anna, run "go tool pprof app.test pprof.cpu"

clean:
	@echo "bash: purging site/rendered and test output"
	@rm -rf site/rendered
	@find test -type d -name rendered -exec rm -rf {} +