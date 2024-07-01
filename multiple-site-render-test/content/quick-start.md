---
date: 2024-04-28
title: Quick Start
toc: true
collections: ["anna"]
---

# Quick Start

---

## Installation

### Installing anna from releases

Run this in the appropriate folder.

If a site directory with the pre-requisite layout template has not been defined, anna proceeds to fetch the default site directory from the anna github repository

```sh
curl -L https://github.com/anna-ssg/anna/releases/download/version-tag/releases-name.tar.gz > anna.tar.gz
tar -xvf anna.tar.gz # unzip the tar file
rm anna.tar.gz # removing the tar file

# here you could add anna to your path if you want and use in in any directory
./anna # runs anna. The instructions are given below
```

### Brew taps for macOS

To install anna on macOS via brew, run the below commands:

```sh
brew tap anna-ssg/anna
brew install anna

# to run anna
anna
```

### Installing anna with go

If you have the Go toolchain installed, run the below command to download and build anna:

```sh
go run github.com/anna-ssg/anna@v2.0.0
```

---

## Usage

### Running anna

1. To run anna, create an `anna.yml` file to configure how sites are rendered and served

#### Sample `anna.yml`

```yml
siteDataPaths:
  - site: site/
  - site-test: site-test/
```

2. Run the render command

- Render all sites

```sh
anna
```

- Render the site located at `site_path`

```sh
anna -r [site_path]
```

### Serving the site with live-reload

1. To serve the site, create an `anna.yml` file to configure how sites are rendered and served

#### Sample `anna.yml`

```yml
siteDataPaths:
  - site: site/
  - site-test: site-test/
```

---

2. Run the serve command

- Serve the site located in `site_path`

```sh
anna -s [site_path]
```

Note: Running `anna -s` without specifying the site_path will throw an error

### Other commands and flags

To view allthe commands and flags available, run the below command:

```sh
  anna -h
```

---
