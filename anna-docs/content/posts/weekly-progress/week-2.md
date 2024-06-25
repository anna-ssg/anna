---
title: Week-2 Progress
date: 2024-03-25
type: post
tags:
  - progress
authors:
  - Adhesh
  - Aditya
  - Anirudh
  - Nathan
---

# Week-2

# Anirudh’s Progress

- Restructured the project
- Improved and integrated the live reload functionality into the ssg
  - Switched to a two goroutine system
    - The main goroutine runs the application and renders pages
    - The second goroutine runs the local web server
  - Eliminated locks and restarting of application on file modification
- Implemented the tagging system
  - Added functionality to tag posts into collections
  - Reverse search for posts of a particular category

---

## Hegde’s Progress

- Implemented chronological feed for posts
- Added ssg flag and frontmatter field to allow working with draft posts. Changed page rendering process to prevent rendering of unnecessary posts.
- Added options in frontmatter and config.yml to integrate javascript based plugins (eg: light mode, code highliting, etc). Users can have seprate plugin options per post.
- Fixed iframe, video and image rendering (CSS)

## Whats Next (Anirudh and Hegde):

- Rebuild the project from ground up and split up the rendering process and Generator struct
- Follow a TDD-based approach during the rebuild

---

# Adhesh’s Progress

- Implemented parallel rendering pipelines.
- Improved parallel rendering and calculation of concurrency factor.
- Cleaned and refactored code to improvise performance.


## What's next?

- Split Parallel rendering pipelines for tags and content files.
- Implement work-stealing.

---

# Nathan’s Progress

- chore/build: Build and deploy anna using Netlify [#11, #48](https://github.com/anna-ssg/anna/pull/48)
- chore/build: switch to Makefile [#49](https://github.com/anna-ssg/anna/pull/49)

## What's next?

- Implement a gui/wizard to configure the `config.yml`

---
