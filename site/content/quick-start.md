---
date: 2024-04-28
title: Quick Start
type: page
toc: true
---

# Quick Start

---

## Installation

### Installing anna from releases

Run this in the appropriate folder. Note that if you don't have a site dir with the pre-requisite layout template; anna proceeds to fetch the default site dir from our GitHub repository

```sh
curl -L https://github.com/anna-ssg/anna/releases/download/version-tag/releases-name.tar.gz > anna.tar.gz
tar -xvf anna.tar.gz # unzip the tar file
rm anna.tar.gz # removing the tar file

# here you could add anna to your path if you want and use in in any directory
./anna # runs anna. The instructions are given below
```

### Brew taps for MacOS

To get anna set up on your mac using brew taps, heres the repo you need to tap off from

```sh
brew tap anna-ssg/anna
brew install anna

# to run anna
anna
```

### Installing anna with go

```sh
go run github.com/anna-ssg/anna@v2.0.0
```

---

## Flags for usage and purpose

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
