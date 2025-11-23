# Anna

[![Test Workflow](https://github.com/anna-ssg/anna/actions/workflows/tests.yml/badge.svg)](https://github.com/anna-ssg/anna/actions/workflows/tests.yml)
[![Netlify Status](https://api.netlify.com/api/v1/badges/09b8bdf3-5931-4295-9fe7-d463d5d06a3f/deploy-status)](https://app.netlify.com/sites/anna-docs/deploys)
[![Go Reference](https://pkg.go.dev/badge/github.com/anna-ssg/v3/anna.svg)](https://pkg.go.dev/github.com/anna-ssg/anna/v3)
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

Pronounced: `/ɐnːɐ/` which means rice 🍚 in Kannada

> This project is a part of the ACM PESU-ECC's yearly [AIEP](https://acmpesuecc.github.io/aiep) program, and is maintained by [Adhesh Athrey](https://github.com/DedLad), [Nathan Paul](https://github.com/polarhive), [Anirudh Sudhir](https://github.com/anirudhsudhir), and [Aditya Hegde](https://github.com/bwaklog)

### Detailed documentation for our SSG can be found [here](https://anna-docs.netlify.app/)

---

## Get Started

```sh
go run github.com/anna-ssg/anna/v3@v3.0.0 -s './'

```

[Quick start](https://anna-docs.netlify.app/quick-start)

<!-- > If you don't have a site directory with the pre-requisite layout template; anna proceeds to fetch the default site dir from our GitHub repository -->

## 📁 Required Site Structure
As of now, Anna is slightly opinionated about how your site must be structured. Currently, Anna expects your site directory to follow the layout shown below:

site/
 ├── content/
 │    ├── index.md
 │    └── posts/
 │         └── sample.md
 │
 ├── layouts/
 │    ├── config.json        # Required
 │    ├── page.html          # Required
 │    ├── posts.html         # Optional (only if you want a posts list page)
 │    └── partials/
 │         └── ...html       # Reusable components such as header/footer
 │
 ├── static/
 │    ├── fonts/
 │    ├── images/
 │    ├── scripts/
 │    ├── style.css
 │    └── plane.jpg          # Example static asset
 │
 └── assets/                 # Optional

## Contributing to Anna

If you have git installed, clone our repository and build against the latest commit

```sh
git clone github.com/anna-ssg/anna; cd anna
go build
```

[Developer Guide](https://anna-docs.netlify.app/developer-guide)
