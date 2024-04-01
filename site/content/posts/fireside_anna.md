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
---

<i>Presented and written by Adhesh, Anirudh, Aditya and Nathan</i>

Building personal blogs from the ground up can be a tedious process. Some of us have had our hands deep in vanilla HTML and CSS, which isn't fun to maintain.
We all want to get to the point and share our thoughts on the web. But, there is a small bump that stops us from doing so

Maintaining your personal site is like working with your own Neovim configuration. Every issue fixed would lead to an entirely unrelated bug.
There is a lot of time spent in fixing things rather than getting productive work done.

A static site generator is an immensely useful application. It can simplify the whole process, allowing you to spend time and energy on quality content.

There are several amazing SSGs out there, like [Hugo](https://gohugo.io/) and [11ty](https://www.11ty.dev/). Building your own SSG is an amazing learning experience. It also motivates one to maintain and improve their personal site.

## Introduction

ACM-PESUECC conducts the ACM Industrial Experience Program(AIEP), an annual program that spans six weeks.
It involves students working as a team to develop an industrial level project.
AIEP intends to give students hands-on experience with real-world projects. It is an excellent opportunity to interact with like-minded individuals.

Our AIEP team consisted of [Adhesh](https://github.com/DedLad), [Aditya](https://github.com/bwaklog), [Nathan](https://github.com/polarhive) and [Anirudh](https://github.com/anirudhsudhir).
Our mentors (.... cool ass senior names) gave us some great ideas for a team of us four freshers. We were puzzled whether to build a distributed postgres clone, or a load balancer. Deep discussions over a week got us to the topic of making blog sites and how tiring it is to work with.

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

Anna is written in Go. We considered writing it in Rust, but that came with a steep learning curve.
Go is a powerful language and has excellent concurrency support, which suited our requirements to build a performant application.

## Whats in the name

Probably the first thing that us four came accross when joining ACM and HSP was the famous Saaru repository.
[Saaru](https://github.com/anirudhRowjee/saaru), one of the projects that inspired our ssg, is derived from a [Kannada](https://en.wikipedia.org/wiki/Kannada) word.
Saaru is a thin lentil soup, usually served with rice.

In Kannada, rice is referred to as 'anna'(ಅನ್ನ) pronounced as <i>/ɐnːɐ/</i>

## Genesis

We began the project in a unique manner, with each of us creating our own prototype of the SSG. This was done to familiarise everyone with the Go toolchain.

The first version of the SSG did just three things.
It rendered markdown(stored in files in a content folder in the project directory) to HTML. This HTML was injected into a layout.html file and served over a local web server.

## Whats unique / what made us want to develop this to a great extent

- Beginner-friendly : Wizard, easy layouts, etc. We want the process of typing out a blog and putting it up on your site to be short and simple.
- Speed : Be fast(hugo,written in Go, is remarkably fast)
- Maintainable : This ssg will be used by us, hence should be easy to fix bugs and add new features

## What is our goal / benchmark

In simple terms, to beat Saaru's render time.
Something you probably never notice while deploying, but it is what pushed us to spend hours on this

## A big rewrite (when we went for a TDD approach)

Starting off this project, we kept adding functionailty without optimisation.
