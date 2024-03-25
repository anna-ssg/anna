# Anna

```text
    ___
   /   |  ____  ____  ____ _
  / /| | / __ \/ __ \/ __ `/
 / ___ |/ / / / / / / /_/ /
/_/  |_/_/ /_/_/ /_/\__,_/

A static site generator in go
```

Inspired by [Hugo](https://gohugo.io) and [Saaru](https://github.com/anirudhRowjee/saaru), this static site generator aims to take performance to the next level with parallel rendering, live reload and so much more, all in Go.

Pronounced: `/ɐnːɐ/` which means rice in Kannada 🍚

This Project is a part of the ACM PESU-ECC's yearly [AIEP](https://acmpesuecc.github.io/aiep) program, and is maintained by [Adhesh Athrey](https://github.com/DedLad), [Nathan Paul](https://github.com/polarhive), [Anirudh Sudhir](https://github.com/anirudhsudhir), and [Aditya Hegde](https://github.com/bwaklog)

---

## Install

Once you have a directory structure, install `anna` using:

```sh
go install github.com/acmpesuecc/anna@latest
```

Or if you have git installed, clone our repository:

```sh
git clone github.com/acmpesuecc/anna --depth=1
cd anna
go run .
```

### The detailed documentation of the SSG can be found [here](https://anna-docs.netlify.app/)\_

---

## Flags

```text
Usage:
  anna [flags]

Flags:
  -a, --addr string     ip address to serve rendered content to (default "8000")
  -d, --draft           renders draft posts
  -h, --help            help for ssg
  -p, --prof            profiles the working code and shows data
  -s, --serve           serve the rendered content
  -v, --validate-html   validate semantic HTML
```
