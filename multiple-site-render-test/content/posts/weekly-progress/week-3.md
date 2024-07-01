---
title: Week-3 Progress
date: 2024-04-01
tags:
  - progress
authors:
  - Adhesh
  - Aditya
  - Anirudh
  - Nathan
collections: ["progress>week-3", "anna", "posts"]
---

# Week-3

# Anirudh and Hegde's Progress

- Split methods previously part of `anna` package among `parser`, `engine` and `helpers` packages
- Refactored `main.go` to only handle flags
- Wrote unit and integration tests for the `parser` and `engine` package
- Split the rendering system to make parallelisation easier by switching to a three method system.
  - Render "anna specific" pages such as sitemap and feed
  - Render "user defined" pages which include all markdown files and posts (This method has been parallelised)
  - Render tags and tag-subpages separately, which could be parallelised in the future
- Wrote a benchmark for `main.go` that times the entire application

## Whats Next (Anirudh and Hegde):

- Improve test coverage for the `engine` package
- Write unit and integration tests for the `cmd` and `helper` packages
- Write unit and integration tests for `main.go`

---

# Adhesh’s Progress

- Re-implemented cobra CLI for the restructed codebase.
- Re-implemented Parallel rendering pipelines for redering tags and content files separately.
- Improved profiling.
- Refactored code to improve performance.
- Worked on content indexing.

## What's next?

- Implement content indexing and site wide content search.
- Improve existing GUI:
  - Add project directory browser.
  - Add Theme browser.
- Implement new flags to provide refined control on resource management.
  - -c / --concurrency to set limit on number of goroutines.
- Implement integration with hosting services to auto host project.

---

# Nathan’s Progress

- Implemented an interactive web based wizard to help a user bootstrap their anna site
  - The intial build wrote the json blob to disk; now it passes it as a POST request to the webserver itself over port `8080` (may conflict)
  - So far it lets you pick a fill metadata, pick a theme and preview your site
  - It auto validates fields using regex as you proceed
  - Also: an animated progressbar and other easter-eggs (confetti??)

## What's next?

- Improve UX by dogfooding and collecting feedback from new users
- Theme triaging and cataloging (basic hugo compatability)
- Bootstrap `site` dir, if the user didn't already have one, using `go-git` and add basic git-submodule support for themes
- Implement go best practices: (auto generate docs, tag releases and known bugs)

---
