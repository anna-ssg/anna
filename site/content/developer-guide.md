---
date: 2024-04-28
title: Developer Guide
toc: false
collections:
  - anna
  - guide
layout: page_2
---

This document explains the repository layout, development commands and how to contribute.

## Repository layout (high level)

- `cmd/anna`        - CLI and server (main command)
- `pkg/engine`      - core rendering engine
- `pkg/parser`      - markdown parsing and helpers
- `site/`           - example site used for docs and dogfooding
- `test/`           - integration and expected output for tests

## Build & run

Build locally:

```bash
go build ./...
```

Run the CLI directly during development:

```bash
go run github.com/anna-ssg/anna/v3 -p ./site -s
```

Run only unit tests:

```bash
go test ./pkg/...
```

Run all tests and race detector:

```bash
go test ./... -v -race
```

## Benchmarks & profiling

Run benchmarks and generate pprof data using the Makefile:

```bash
make bench
# results and pprof files are in the profiles/ or test output directories
```

While serving you can view pprof endpoints at `http://localhost:8000/debug/pprof` (see `cmd/anna`).

## Makefile targets

```
Targets:
  build  : Build anna and render the site
  serve  : Build anna, render and serve the site with live reload
  tests  : Run all tests
  bench  : Run the benchmark and generate pprof files
  clean  : Remove the rendered site directory and test output
```

## Contribution workflow

1. Fork the repository and create a feature branch (eg. `feature/foo`).
2. Write code, run `go fmt ./...` and add/update tests.
3. Run `go test ./... -v` and verify all tests pass.
4. Push your branch and open a Pull Request with a clear description and changelog.

### Development tips

- Use `go vet` and `go test -race` while developing for correctness.
- Keep changes small and isolated; add tests for new behavior.
- Update `site/` when adding features that affect rendering or layouts so the docs can dogfood changes.

---