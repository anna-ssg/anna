all: build

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "")

LDFLAGS := \
	-X main.Version=$(VERSION) \
	-X main.FullCommitHash=$(COMMIT)

build:
	@echo "anna: building $(VERSION) ($(COMMIT))"
	@go build -ldflags "$(LDFLAGS)"
	@./anna

serve:
	@echo "anna: serving site"
	@go build -ldflags "$(LDFLAGS)"
	@./anna -s -p site

tests:
	@echo "anna: running all tests"
	@go test ./...

bench:
	@echo "anna: running benchmark"
	@go test -bench . -benchmem -cpuprofile pprof.cpu
	@# to profile anna, run "go tool pprof app.test pprof.cpu"

clean:
	@echo "bash: purging site/rendered and test output"
	@rm -rf site/rendered
	@find test -type d -name rendered -exec rm -rf {} +