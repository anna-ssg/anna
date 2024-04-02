---
title: Presenting anna at fireside
date: 2024-04-03
type: post
description: Building a static site generator in Go
draft: false
tags:
  - acm
  - hsp
  - go
  - tech
  - talk
  - aiep
---

*Presented and written by Adhesh, Anirudh, Aditya and Nathan*

Building personal blogs from the ground up can be a tedious process. Some of us have had our hands deep in vanilla HTML and CSS, which isn't fun to maintain.We all want to get to the point and share our thoughts on the web. But, there is a small bump that stops us from doing so

Maintaining your personal site is like working with your own Neovim configuration. Every issue fixed would lead to an entirely unrelated bug.There is a lot of time spent in fixing things rather than getting productive work done.

A static site generator is an immensely useful application. It can simplify the whole process, allowing you to spend time and energy on quality content.

There are several amazing SSGs out there, like [Hugo](https://gohugo.io/) and [11ty](https://www.11ty.dev/). Building your own SSG is an amazing learning experience. It also motivates one to maintain and improve their personal site.

## Introduction

ACM-PESU ECC conducts the ACM Industrial Experience Program(AIEP), an annual program that spans six weeks. 

> It involves students working as a team to develop an industrial level project. AIEP intends to give students hands-on experience with real-world projects. It is an excellent opportunity to interact with like-minded individuals.

Our AIEP team consisted of [Adhesh](https://github.com/DedLad), [Aditya](https://github.com/bwaklog), [Nathan](https://github.com/polarhive) and Anirudh. 

Our mentors (cool ass senior names!) gave us some great ideas for a team of us four freshers. We were puzzled whether to build a distributed postgres clone, or a load balancer. Deep discussions over a week got us to the topic of making blog sites and how tiring it is to work with which only gets worse as you write more and more content for your internet home.

This and inspirations from [Saaru](https://github.com/anirudhRowjee/saaru) and [Sapling](https://github.com/NavinShrinivas/sapling) pushed us to tackle this problem with our own SSG.

```text
    ___
   /   |  ____  ____  ____ _
  / /| | / __ \/ __ \/ __ `/
 / ___ |/ / / / / / / /_/ /
/_/  |_/_/ /_/_/ /_/\__,_/

A static site generator in go
```

## The small but big decision

Anna is written in Go. We considered writing it in Rust, but that came with a steep learning curve. Go is a powerful language and has excellent concurrency support, which suited our requirements to build a performant application.

## What's in the name

Probably the first thing that us four came across when joining ACM and HSP was the famous Saaru repository. [Saaru](https://github.com/anirudhRowjee/saaru), one of the projects that inspired our ssg, is derived from a [Kannada](https://en.wikipedia.org/wiki/Kannada) word. Saaru is a thin lentil soup, usually served with rice.

In Kannada, rice is referred to as 'anna'( ಅನ್ನ) pronounced <i>/ɐnːɐ/</i>

---
# Genesis

We began the project in a unique manner, with each of us creating our own prototype of the SSG. This was done to familiarise everyone with the Go toolchain.

The first version of the SSG did just three things. It rendered markdown (stored in files in a content folder in the project directory) to HTML. This HTML was injected into a layout.html file and served over a local web server.

## What made us want to develop this to a great extent

- Beginner-friendly: Wizard, easy ready to use layouts, etc. We want the process of typing out a blog and   putting it up on your site to be short and simple.
- Speed: Be Fast (hugo – written in Go, is remarkably fast)
- Maintainable: This ssg will be used by us, hence it should be easy to fix bugs and add new features
- Learning curve: none of us have really shipped a production ready application and since AIEP is all about making industry ready project, we chose to use Go, so we could spend more *writing* code and not worrying about our toolchain or dependency hell.
- Owning a piece of the internet: Aditya and Anirudh registered their own domain names. Now their anna sites live on hegde.live and sudhir.live.

## goals / benchmarks?

In simple terms, to beat Saaru's render time. Something you probably never notice while deploying, but it is what pushed us to spend hours on this.

Adhesh was pretty excited to pick up go and implement profiling, shaving milliseconds off of builds as he implemented parallel rendering using goroutines.

## A big rewrite (when we went for a TDD approach)

Starting off this project, we kept adding functionality without optimisation. We didn’t have a proper structure, PRs would keep breaking features and overwriting other functions. We proceeded to restructure our ssg into:

> modules previously part of `cmd/anna/utils.go` and `cmd/anna/main.go` are to be split between `pkg/parsers/`, `pkg/engine/` and `pkg/helper`

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

### New proposed rendering system

- Currently there are two separate types of files that are rendered. One set are user defined files whicha are for example docs.md, timeline.md. These are specific to a user.
- The second set of files that are rendered are files that belong to `tags.html`, `sub-tags.html` and `posts.html`
- Splitting the rendering system to make parallelisation easier. Now the generator/engine has a method to render "anna specific" pages and another method to render "user defined" pages which include all the user pages and posts

---
## To search or not to search?? 🤔

We were wondering if we’d need a search function on our site, since google.com and any other search engine indexes our site anyway. If we needed to implement search, we had a constraint: we could not use an API, it had to be static and local to be user-friendlyl. Aditya and Anirudh implemented a json func that uses  deepdatamerge to index posts on our site. // to be filled

## Wizard
// to be filled
// add pix
Nathan proceeded to work on a GUI; a web-based wizard that let's a new user setup anna along with a couple of easter eggs along the way 🍚

The wizard lets a user pick a theme, enter your name, pick navbar elements, and validates fields using regex checks so you don’t need to worry about relative paths in baseURLs, canonical links, and sitemaps. After successfully completing the setup the wizard launches a live preview of your site in a new tab.


