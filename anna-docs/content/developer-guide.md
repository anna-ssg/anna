---
date: 2024-04-28
title: Developer Guide
type: page
toc: true
collections:
  - anna
  - guide
layout: page_2
---

# Developer Guide

---

## Contributing to Anna

Detailed documentation for our SSG can be found: [here](https://anna-docs.netlify.app/)

If you have git installed, clone our repository and build against the latest commit

```sh
git clone github.com/anna-ssg/anna; cd anna
go build
```

---

### Profiling

The live profile data of the application can be viewed during live reload by navigating to `http://localhost:8000/debug/pprof`

---

## Makefile

The Makefile contains various commands to aid development

```text
Usage:
  make [target]

Targets:
  build: Build anna and render the site
  serve: Build anna, render and serve the site with live reload
  tests: Run all tests
  bench: Run the benchmark and generate pprof files
  clean: Remove the rendered site directory and test output
```
