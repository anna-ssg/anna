---
date: 2024-04-28
title: Quick Start
toc: true
collections: ["anna"]
---

# Quick Start

---

## Installation

### Brew taps for macOS

To install anna on macOS via brew, run the below commands:

```sh
brew tap anna-ssg/anna
brew install anna

# to run anna
anna
```

### Installing anna from releases

```sh
curl -L https://github.com/anna-ssg/anna/releases/download/version-tag/releases-name.tar.gz > anna.tar.gz
tar -xvf anna.tar.gz # unzip the tar file
rm anna.tar.gz # removing the tar file

# here you could add anna to your path if you want and use in in any directory
./anna # runs anna. The instructions are given below
```

### Installing anna with go

If you have the Go toolchain installed, run the below command to download and build anna:

```sh
go run github.com/anna-ssg/anna@v3.0.0
```

---

## Usage

### Running anna

```sh
anna
```

- Render the site located at `site_path`

```sh
anna -p [site_path]
```

### Serving the site with live reload to watch for file updates

```sh
anna -s
```

- Serve the site located in `site_path`

```sh
anna -s -p [site_path]
```

### Other commands and flags

To view all the commands and flags available, run the below command:

```sh
  anna -h
```

---
