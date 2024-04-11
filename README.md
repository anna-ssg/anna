# Anna

[![Test Workflow](https://github.com/acmpesuecc/anna/actions/workflows/tests.yml/badge.svg)](https://github.com/acmpesuecc/anna/actions/workflows/tests.yml)
[![Netlify Status](https://api.netlify.com/api/v1/badges/09b8bdf3-5931-4295-9fe7-d463d5d06a3f/deploy-status)](https://app.netlify.com/sites/anna-docs/deploys)
[![Go Reference](https://pkg.go.dev/badge/github.com/acmpesuecc/anna.svg)](https://pkg.go.dev/github.com/acmpesuecc/anna)
[![GitHub Repo Stars](https://img.shields.io/github/stars/acmpesuecc/Anna?style=flat-square&label=Stars&color=lightgreen&logo=github)](https://github.com/acmpesuecc/Anna)

```text
    ___
   /   |  ____  ____  ____ _
  / /| | / __ \/ __ \/ __ `/
 / ___ |/ / / / / / / /_/ /
/_/  |_/_/ /_/_/ /_/\__,_/

A static site generator in go
```

Inspired by [Hugo](https://gohugo.io) and [Saaru](https://github.com/anirudhRowjee/saaru), this static site generator aims to take performance to the next level with parallel rendering, live reload and so much more, all in Go.

Pronounced: `/ɐnːɐ/` which means rice 🍚 in Kannada

> This project is a part of the ACM PESU-ECC's yearly [AIEP](https://acmpesuecc.github.io/aiep) program, and is maintained by [Adhesh Athrey](https://github.com/DedLad), [Nathan Paul](https://github.com/polarhive), [Anirudh Sudhir](https://github.com/anirudhsudhir), and [Aditya Hegde](https://github.com/bwaklog)

---
## Get Started

```sh
go run github.com/acmpesuecc/anna@v1.0.0-alpha
```
> If you don't have a site dir with the pre-requisite layout template; anna proceeds to fetch the default site dir from our GitHub repository

## Contributing to Anna

Detailed documentation for our SSG can be found: [here](https://anna-docs.netlify.app/)

If you have git installed, clone our repository and build against the latest commit

```sh
git clone github.com/acmpesuecc/anna; cd anna 
go build
```
```text
Usage:
  anna [flags]

Flags:
  -a, --addr string   ip address to serve rendered content to (default "8000")
  -d, --draft         renders draft posts
  -h, --help          help for anna
  -l, --layout        validates html layouts
  -p, --prof          enable profiling
  -s, --serve         serve the rendered content
  -v, --version       prints current version number
  -w, --webconsole    wizard to setup anna
```
