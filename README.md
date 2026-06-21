# Anna

[![Test Workflow](https://github.com/anna-ssg/anna/actions/workflows/tests.yml/badge.svg)](https://github.com/anna-ssg/anna/actions/workflows/tests.yml)
[![Netlify Status](https://api.netlify.com/api/v1/badges/09b8bdf3-5931-4295-9fe7-d463d5d06a3f/deploy-status)](https://app.netlify.com/sites/anna-docs/deploys)
[![Go Reference](https://pkg.go.dev/badge/github.com/anna-ssg/v3/anna.svg)](https://pkg.go.dev/github.com/anna-ssg/anna/v4)
[![GitHub Repo Stars](https://img.shields.io/github/stars/anna-ssg/Anna?style=flat-square&label=Stars&color=lightgreen&logo=github)](https://github.com/anna-ssg/anna)

```text
    ___
   /   |  ____  ____  ____ _
  / /| | / __ \/ __ \/ __ `/
 / ___ |/ / / / / / / /_/ /
/_/  |_/_/ /_/_/ /_/\__,_/

A static site generator in go
```

Inspired by [Hugo](https://gohugo.io) and [Saaru](https://github.com/anirudhRowjee/saaru), this static site generator aims to take performance to the next level with parallel rendering, live reload and so much more, all in Go.

> Pronounced: `/ɐnːɐ/` which means rice 🍚 in Kannada

---

## Get Started

> To setup your site, follow the quick-start [guide](https://anna-docs.netlify.app/quick-start)

## Contributing to Anna

If you have git installed, clone our repository and build against the latest commit

```sh
git clone github.com/anna-ssg/anna
cd anna
make build
```

`make build` automatically wires up a **pre-commit hook** (tracked in [`.github/githooks/`](.github/githooks/)) that runs, in order: `gofmt` check, lint (`go vet` + [`golangci-lint`](https://golangci-lint.run/welcome/install/) if installed), audit ([`govulncheck`](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) if installed), the test suite, `golangci-lint` and `govulncheck` are optional; the hook skips those steps with a hint if they aren't installed.

You can also install the hook explicitly, without doing a full build, by running:

```sh
make install-hooks
```

If you ever need to skip it (e.g. a WIP commit), use `git commit --no-verify`.

### Developer Guide

Detailed documentation for developers can be found [here](https://anna-docs.netlify.app/developer-guide)

---

### History

> *This project was a part of the ACM PESU-ECC's yearly [AIEP](https://acmpesuecc.github.io/aiep) program, and is maintained by [Adhesh Athrey](https://github.com/DedLad), [Nathan Paul](https://github.com/polarhive), [Anirudh Sudhir](https://github.com/anirudhsudhir), and [Aditya Hegde](https://github.com/bwaklog)*