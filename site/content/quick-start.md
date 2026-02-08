---
date: 2024-04-28
title: Quick Start
toc: false
collections: ["anna"]
---

## Installation

> Download a prebuilt binary for your platform from the latest release:

### MacOS

```sh
brew install anna-ssg/anna/anna
```


### GNU/Linux (tar.gz):

> 1. Linux (x86_64): [anna_Linux_x86_64.tar.gz](https://github.com/anna-ssg/anna/releases/latest/download/anna_Linux_x86_64.tar.gz)
> 2. `tar -xzf anna_Linux_x86_64.tar.gz`
> 3. `./anna`

### Windows 10/11 (x86_64):

> 1. Download [anna_Windows_x86_64.zip](https://github.com/anna-ssg/anna/releases/latest/download/anna_Windows_x86_64.zip)
> 2. Unzip `anna_Windows_*.zip`
> 3. Run `anna.exe` from the extracted folder in a terminal

---

## Bootstrap (create a `site/` directory)

If you don't already have a `site/` directory, Anna can initialize one for you with a default layout.

```sh
# will ask to download the default site layout when config is missing
./anna 
```

---

## Usage

### Basic render

Render the site found at `site/` (default):

```sh
anna
```

Specify a path:

```sh
anna -p ./site
```

### Serve with live reload

```sh
anna -s
```

### Version and debug

Show the version (includes embedded commit when present):

```sh
anna -v
```

Show usage and all flags:

```sh
anna -h
```
