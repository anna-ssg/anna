---
title: week-0
date: 2024-03-02
type: post
draft: false
tags:
  - progress
---

# week-0

# Anirudh’s Progress

- Use the 'Goldmark' library to generate HTML from markdown content
- Create web server to serve rendered content
- Create an initial working prototype of front matter.
- Render the content/ and static/ directory recursively
- Create a reusable system of partials to build page layouts
- Parse the titles of posts from config.yml, later changed

## What's next?

- tofill

---

# Adhesh’s Progress

- Migrated to [Cobra](https://cobra.dev) for better **CLI integration**.
- Implemented important flags for serving a local **HTTP server** (---serve), and to explicitly mention the **port** to use (---addr).
- Implemented **Real-Time Directory/File change** watcher using [fsnotify](https://pkg.go.dev/github.com/fsnotify/fsnotify), and Hot-Reload system for reserving updated files instantaneously.
- Cleaned and sanitised some parts of the codebase.
- Included early optimizations to some functions of the codebase using goroutines and sync operations.

## What's next?

- Making **main.go** the entry point for the code.
- Fix issues with the watcher, and clean up the goroutine issues.
- Try studying about **parallelizing code functions**, and implementing it.
- Add a **developer** mode/flag for profiling the performance of the application (---dev)
- Look into [Huge theme](https://themes.gohugo.io) compatibility.
- Try integrating **JavaScript** in templates.

---

# Hegde’s Progress

- Switched to automatic filename parsing
- Implementing Draft Posts
- Complete CSS styling
- Cleaned up unnecessary post rendering

## Whats Next:

- JS inection as plugins into pages and individual posts
- Draft post rendering
- Chronological Feed for posts

# Nathan’s Progress

- Setup CI using GitHub actions which builds and deploys ssg to [gh-pages](https://ssg-test-org.github.io)
- Setup `robots.txt` to (currently set to prevent indexing until we reach a v1.0 release)
- Implemented `sitemap.xml` to tell search engines how our site is structured
- Fixed `baseURL` it now, uses absolute paths from root `/` when loading stylesheets

## What's next?

- SEO optimization (verify semantic HTML generation)
- OGP tags in page headers

---
