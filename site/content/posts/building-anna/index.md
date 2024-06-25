---
title: Building anna
date: 2024-04-04
type: post
draft: false
description:
  This page contains a post about anna, a static site generator written in Go. This team
  project was built as part of AIEP 2024
tags:
  - acm
  - hsp
  - go
  - tech
  - talk
  - aiep
authors:
  - Adhesh
  - Aditya
  - Nathan
  - Anirudh
---

There are several amazing SSGs out there, like [Hugo](https://gohugo.io/) and
[11ty](https://www.11ty.dev/). Building your own SSG is an amazing learning
experience. It also motivates one to maintain and improve their personal site.

> Presented and written by Adhesh, Anirudh, Aditya and Nathan

Building personal blogs from the ground up can be a _tedious process_. Some of us
have had our hands deep in vanilla HTML and CSS, which isn't fun to maintain. We
all want to get to the point and share our thoughts on the web. But, there is a
small bump that stops us from doing so.

Maintaining your personal site is like working with your own Neovim
configuration. Every issue fixed would lead to an entirely unrelated bug. There
is a lot of time spent fixing things rather than getting productive work done.

> A static site generator is an immensely useful application

![Lighthouse scores of the anna-docs page](images/lighthouse.png)

It can simplify the whole process: allowing you to spend time and energy
on quality content. Keep reading to find out how we designed anna `v1.0.0`

---

## Introduction

ACM-PESU ECC conducts the ACM Industrial Experience Program (AIEP), an annual program that spans six weeks.

> It involves students working as a team to develop an industrial-level
> project. AIEP intends to give students hands-on experience with real-world
> projects. It is an excellent opportunity to interact with like-minded
> individuals.

Our AIEP team consisted of [Adhesh](https://github.com/DedLad), [Aditya](https://github.com/bwaklog),
[Nathan](https://github.com/polarhive), and [Anirudh](https://github.com/anirudhsudhir).

Our mentors (cool ass senior names!) gave us some great ideas for a team of us
four freshers.
We were puzzled whether to build a distributed Postgres clone or a load balancer.

Deep discussions over a week got us to the topic of making
blog sites and how tiring it is to work with, which only gets worse as you
write more and more content for your internet home.

This and inspirations from [Saaru](https://github.com/anirudhRowjee/saaru) and
[Sapling](https://github.com/NavinShrinivas/sapling) pushed us to tackle this
problem with our own SSG.

```text
    ___
   /   |  ____  ____  ____ _
  / /| | / __ \/ __ \/ __ `/
 / ___ |/ / / / / / / /_/ /
/_/  |_/_/ /_/_/ /_/\__,_/

A static site generator in Go

```

## The small but big decision!

Anna is written in [Go](https://go.dev). We considered writing it in Rust, but
that came with a steep learning curve.
Go is a powerful language and has excellent concurrency support, which suited our requirements to build a performant application.

### What's in a name?

Probably the first thing that the four of us came across when joining ACM and HSP was the famous Saaru repository.
[Saaru](https://github.com/anirudhRowjee/saaru),
one of the projects that inspired our ssg,
is derived from a [Kannada](https://en.wikipedia.org/wiki/Kannada) word.
Saaru is a thin lentil soup, usually served with rice.

> In Kannada, rice is referred to as 'anna'(ಅನ್ನ) pronounced <i>/ɐnːɐ/</i>

This was just a playful stunt that we engaged in. We planned on beating Saaru at
site render times, optimizing at runtime.

---

## Genesis

We began the project in a unique manner, with each of us creating our own
prototype of the SSG. This was done to familiarize everyone with the Go
toolchain.

The first version of the SSG did just three things. It rendered markdown
(stored in files in a content folder in the project directory) to HTML. This
HTML was injected into a layout.html file and served over a local web server.
Later, we implemented a front matter YAML parser to retrieve page metadata

---

## What made us develop this to a great extent?

- Beginner-friendly: An easy setup wizard, easy and ready to use layouts, and themes. We want the
  process of typing out a blog and putting it up on your site to be short and
  simple.
- Speed: Be fast (hugo – written in Go, is remarkably fast)
- Maintainable: This ssg will be used by us, hence it should be easy to fix
  bugs and add new features
- Learning curve: None of us have really shipped a production ready
  application. Since AIEP is all about making industry-ready projects, we chose
  to use go: so we could spend more **\*writing** code\* and not worrying about our
  toolchain or escaping dependency hell.
- Owning a piece of the internet: Aditya and Anirudh registered their own
  domain names. Now their anna sites live on [hegde.live] and [sudhir.live]

---

## Benchmarks! Can anna lift??

In simple terms, to beat Saaru's render time (P.S. we did!). Something you
probably never notice while deploying, but it is what pushed us to spend hours
on this.

Adhesh was pretty excited to pick up Go and implement profiling, shaving
milliseconds off build times, when he implemented parallel rendering using
goroutines.

### We cook! 🍳

Here are some screenshots out of our group chats, that demonstrate build times, profiling et-al when having thousands of markdown files or in this case
just copy-pasting a single markdown file en-mass!

![Hyperfine benchmarks comparing the render times of anna, Saaru and 11ty](images/bench.png)
Hyperfine benchmarks comparing the render times of anna, Saaru and 11ty

> After about 2 weeks of training (_ahem_) coding, we had a (merge) bringing parallel rendering and profiling to the table

---

## Profiling (WIP)

For those who love to get technical, anna has a profiling flag.
This flag prints the render profile data to stdout.
It also writes a cpu.prof and mem.prof file, which can be analysed and visualised using [pprof](https://pkg.go.dev/runtime/pprof)

Here,s the CPU profile generated while rendering this site.
This is an SVG showing how much time each function call takes, the number of times it ran in a given test sample and various other useful information.

![CPU profile of an anna render generated using pprof](images/cpu_prof.svg)

<!-- ![CPU profile of an anna render generated using pprof](images/cpu_prof.svg) -->

You may wanna zoom-in about 3-4x times to get to see how our ssg works

---

## Live Reload

### Prototype

The initial prototype built by Adhesh consisted of a multi-goroutine system.
A new goroutine would be spawned to walk the required directories.
If the current path being walked was a file, the path would be passed to another function along with its current modification time.

The previous mod time of the file would then be retrieved from a map holding the mod times of all the files:

- If the given file was freshly created, its modification time would be added to the map.
- If there was no change in the mod time, no changes would be made.
- If there was a change between the current and previous mod times, another function would be called.

The new function checks if a child process is running:

- For the first render, when a process has not been created, a new process is created that runs anna ("go run main.go --serve")
- For successive renders, the existing process is killed and a new process is spawned once again that runs anna.

This prototype was not very efficient as it created and killed processes for every change.
It had multiple goroutines attempting to walk the directories at the same time.
It also used multiple mutual exclusion locks to prevent data races.
Integrating this into the project also proved to be challenging.

### Improved version

The live reload feature was improved by Anirudh.
The updated version utilised two goroutines.

The main goroutine used the earlier file walker, with one important change: it sequentially traversed the directory without spawning new goroutines.
For any modification to a file in the current traversal, a vanilla render of the entire site would be performed.
The goroutine would then sleep for a specified duration (currently 1 second) before attempting the next directory traversal.

The secondary goroutine ran a local web server that served the rendered/ directory.

This eliminated all locks and avoided the creation and destruction of any child processes.

---

## A big rewrite (when we went for a TDD approach) (WIP)

Starting off this project, we kept adding functionality without optimization.
We didn’t have a proper structure; PRs would keep breaking features and overwriting functions written by team-mates.

Anirudh and Aditya dived into the codebase to rebuild the application following a TDD-approach.

### A new proposed rendering system

#### main.go

Firstly, we refactored main.go to only handle flags.
The rest of the program logic was moved to other packages.
A benchmark for main.go was also written to time the entire application.

#### pkg/ modules

Modules previously part of `cmd/anna/utils.go` and `cmd/anna/main.go` were to be split between `pkg/parsers/`, `pkg/engine/` and `pkg/helper`

```text
pkg
├─── helpers
│   ├─── helpers.go
│   └─── helper_test.go
├─── engine
│   ├─── anna_engine.go
│   ├─── anna_engine_test.go
│   ├─── engine.go
│   ├─── engine_test.go
│   ├─── user_engine.go
│   ├─── user_engine_test.go
│   └─── engine_integration_test.go
└─── parsers
	├── parser.go
	├── parser_test.go
	└── parser_integration_test.go
```

> - Wrote unit and integration tests for the parser and engine package

#### Splitting the render pipeline

Currently, there are three kinds of files that have to be rendered.
One set includes user-defined files such as `index.md`, `docs.md` and various posts. These are specific to a user.

The second set of files that are rendered include `tags.html` and `posts.html`, which are present on every site rendered by anna.

The third set of files included the tag-sub pages.
For every tag, there would be a corresponding sub-page containing all of the posts with the same tag.

> - Split the rendering system to make parallelisation easier by switching to a three method system.

> - Render "user defined" pages which include all markdown files and posts (This method has been parallelised, Render tags and tag-sub pages separately, which could be parallelised in the future

---

## Tags (WIP)

You can tag posts by hand, at the start of each markdown file and you get a
nice sub-page on your site so readers can discover similar content or browser
by category.

- Organizing posts into collections based on tags
- Reverse search for posts of a certain category

---

## To search or not to search? 🤔

> That is the question > Is our _static site_ becoming dynamic and at what cost?

We were wondering if we’d need a search function on our site since Google and
any other web-crawler index our site anyway.
If we needed to implement it, we had a constraint: do not use an API.
It had to be static and local to be user-friendly to work with.
Aditya and Anirudh implemented a JSON index generator that uses "Deep Data Merge" to index posts on our site.

This index is generated during the site render and functions without any lag or noticeable delay when searching across posts.
We mean to re-write it using WASM if necessary.

Here's a gif demonstrating search

![Demonstration of the search feature in anna](images/search.gif)

## JS integration as plugins

Aditya added a field to our frontmatter which lets you pick and add certain JS
based snippets to your site.
This way, you get to add `highlight.js` support, analytics scripts and donation page widgets; that you can source from the `static/scripts` folder and toggle as needed per-markdown page.

### Raw HTML

What if you'd want to add a contact form to your site? or embed YouTube videos or iframes of your choosing?

Anna let's us do that! Although, the point of a static site generator is to
quickly get to writing and focusing on the content.
You can still embed js elements and iframe as needed to showcase any interesting YouTube videos or to just rickroll people!

---

## Wizard

Nathan proceeded to work on a GUI; a web-based wizard that let's a new user
setup anna along with a couple of easter eggs along the way 🍚

The wizard lets a user pick a theme, enter your name, pick navbar elements, and
validates fields using regex checks so you don’t need to worry about relative
paths in baseURLs, canonical links, and sitemaps. After successfully completing
the setup, the wizard launches a live preview of your site in a new tab.

![Demonstration of the GUI wizard in anna](images/wizard.gif)

---

### changelog: showcasing important additions, which are yet to be added to this blog

Nathan:

- feat: implement sitemap.xml by @polarhive in #17
- feat: ogp tags and atom feed by @polarhive in #33
- feat: bootstrap site and import stylesheets by @polarhive in #73

Aditya:

- feat: implement drafts by @bwaklog in #9
- feat: chronological feed, js plugins (eg: light.js, prism.js) by @bwaklog in #32

---

## Feedback? / Learnings

We are at week: 4/6 and have a lot of things in store and bugs to squash!

> Feel free to ask any questions / send feature requests you'd like to see?

This blog post misses out of many not-so well documented features and learnings that
we got during midnight calls and the patches we kept sending each other fixing trivial but
interesting issues.
Have a look at our [GitHub](https://github.com/anna-ssg/anna/issues), for more

---

Today [anna](https://github.com/anna-ssg/anna/releases/latest) is tagged at v1.0.0 and we use it on our personal sites:
[hegde.live](https://hegde.live) // [sudhir.live](https://sudhir.live) // [polarhive.net](https://polarhive.net)

---

01100001 01101110 01101110 01100001
