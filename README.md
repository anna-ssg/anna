# Anna

> _`/ɐnːɐ/`_ is a Kannada word for 'rice' 🍚

Inspired by [Hugo](https://gohugo.io) and [Saaru](https://github.com/anirudhRowjee/saaru). Anna is a lightning fast static site generator written in Go, designed for simplicity and ease of use. With a focus on performance and minimal configuration, Anna lets you to create beautiful static websites with ease!

```text
    ___
   /   |  ____  ____  ____ _
  / /| | / __ \/ __ \/ __ `/
 / ___ |/ / / / / / / /_/ /
/_/  |_/_/ /_/_/ /_/\__,_/

A static site generator in go
```

---

## Get Started!

[![GitHub Repo Stars](https://img.shields.io/github/stars/anna-ssg/Anna?style=flat-square&label=Stars&color=lightgreen&logo=github)](https://github.com/anna-ssg/anna)
[![Anna Docs](https://img.shields.io/badge/anna-docs-0d7ebf)](https://anna-docs.netlify.app)
[![Test Workflow](https://github.com/anna-ssg/anna/actions/workflows/tests.yml/badge.svg)](https://github.com/anna-ssg/anna/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/anna-ssg/v3/anna.svg)](https://pkg.go.dev/github.com/anna-ssg/anna/v4)

> To setup your site, follow our [quick-start](https://anna-docs.netlify.app/quick-start) guide to download Anna

### Deploy Your Site

Anna ships as a single static binary, so most hosts deploy it the same way: fetch the `anna` release binary, run it to generate `site/rendered`, then publish that folder. No Go toolchain needed on the host's end.

#### Netlify

1. Push a repo containing your `site/` directory
2. Copy [`deploy.sh`](deploy.sh) and [`netlify.toml`](netlify.toml) into it
3. Connect the repo on [Netlify](https://app.netlify.com) — build command and publish dir are picked up automatically from your `netlify.toml` file

#### Cloudflare Pages

1. Connect your repo on [Cloudflare Pages](https://pages.cloudflare.com)
2. Build command: `bash deploy.sh`
3. Build output directory: `site/rendered`

#### GitHub Pages

1. Settings → Pages → Source: **GitHub Actions**
2. Look at the **Build and Deploy** workflow from our Actions tab
3. `deploy.sh` fetches a `Linux_x86_64` release binary, which matches the default Linux runners on all three platforms above.

---
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

### History

> *This project was a part of the ACM PESU-ECC's yearly [AIEP](https://pesuecc.acm.org/aiep/2024) program, and is maintained by [Adhesh Athrey](https://github.com/DedLad), [Nathan Paul](https://github.com/polarhive), [Anirudh Sudhir](https://github.com/anirudhsudhir), and [Aditya Hegde](https://github.com/bwaklog)*
