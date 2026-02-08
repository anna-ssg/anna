all:
	@make build

COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "")
COMMIT_SHORT ?= $(shell echo $(COMMIT) | cut -c1-7)
LDFLAGS ?= -X main.CommitHash=$(COMMIT)
build:
	@echo "anna: building with commit $(COMMIT_SHORT)"
	@go build -ldflags "$(LDFLAGS)"
	@./anna

serve:
	@echo "anna: serving site"
	@go build
	@./anna -s "site/"
tests:
	@echo "anna: running all tests"
	@go test ./...
bench:
	@echo "anna: running benchmark"
	@go test -bench . -benchmem -cpuprofile pprof.cpu
	@#to profile anna, run "go tool pprof app.test pprof.cpu"
clean:
	@echo "bash: purging site/rendered and test output"
	@rm -rf site/rendered
	@cd test/
	@rm -rf `find . -type d -name rendered`
	@cd ../

# Build with embedded commit hash (set COMMIT env var or default to current git HEAD)
