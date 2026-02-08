all:
	@make build

# Build with embedded commit hash (set COMMIT env var or default to current git HEAD)
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "")
LDFLAGS ?= -X main.CommitHash=$(COMMIT)
build:
	@echo "anna: building with commit $(COMMIT)"
	@go build -ldflags "$(LDFLAGS)"
	@./anna

serve:
	@echo "anna: serving site"
	@go build
	@./anna -s -p "site"
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
