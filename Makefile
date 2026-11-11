all: build

.PHONY: all install-hooks build tests tests-coverage fmt fmt-check lint govulncheck bench clean

install-hooks:
	@git rev-parse --is-inside-work-tree >/dev/null 2>&1 && \
		git config core.hooksPath .github/githooks && \
		echo "pre-commit hook installed"

build:
	go build .

tests:
	go test -v -failfast ./...

tests-coverage:
	go test -v -coverpkg=./... -coverprofile=coverage.out ./...

fmt:
	gofmt -l -w .
	gofmt -l .

fmt-check:
	gofmt -l .

lint:
	golangci-lint run ./...; \

govulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	
bench:
	go test -bench . -benchmem -cpuprofile pprof.cpu
	@printf "\nTo profile anna, run \"go tool pprof app.test pprof.cpu\"\n"

clean:
	rm -rf site/rendered
