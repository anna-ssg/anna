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

## Deploy Your Site

Anna ships as a single static binary, so most hosts deploy it the same way: fetch the `anna` release binary, run it to generate `site/rendered`, then publish that folder. No Go toolchain needed on the host's end.

### Netlify

1. Push a repo containing your `site/` directory
2. Copy [`deploy.sh`](deploy.sh) and [`netlify.toml`](netlify.toml) into it
3. Connect the repo on [Netlify](https://app.netlify.com) — build command and publish dir are picked up automatically from your `netlify.toml` file

### Cloudflare Pages

1. Connect your repo on [Cloudflare Pages](https://pages.cloudflare.com)
2. Build command: `bash deploy.sh`
3. Build output directory: `site/rendered`

### GitHub Pages

1. Settings → Pages → Source: **GitHub Actions**
2. Look at the **Build and Deploy** workflow from our Actions tab
3. `deploy.sh` fetches a `Linux_x86_64` release binary, which matches the default Linux runners on all three platforms above.