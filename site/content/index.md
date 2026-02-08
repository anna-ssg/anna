---
date: 2024-02-24
title: Henlo!
previewimage: "/static/images/anna.png"
description: "Homepage+Docs for Anna"
layout: page
---

Anna is a lightning fast static site generator written in Go, designed for simplicity and ease of use. With a focus on performance and minimal configuration, [Anna](https://github.com/anna-ssg/anna) allows you to create beautiful static websites with ease.

## Crafted with Anna

<div id="embed-gallery"></div>

<script>
const urls = [
  "https://hsp-ec.xyz",
  "https://anirudhsudhir.com",
  "https://adihegde.com",
  "https://polarhive.net"
];

const gallery = document.getElementById('embed-gallery');
gallery.className = 'embed-gallery';

urls.forEach((url, index) => {
  const domain = new URL(url).hostname;
  const item = document.createElement('div');
  item.className = 'embed-item';
  item.style.animationDelay = `${index * 0.5}s`;
  const frame = document.createElement('div');
  frame.className = 'embed-frame';
  const iframe = document.createElement('iframe');
  iframe.src = url;
  iframe.title = `Live site preview — ${domain}`;
  iframe.loading = 'lazy';
  iframe.sandbox = 'allow-forms allow-scripts allow-same-origin allow-popups';
  iframe.referrerPolicy = 'no-referrer';
  frame.appendChild(iframe);
  item.appendChild(frame);
  gallery.appendChild(item);
});
</script>

---

# [Quick Start Guide](./quick-start)

Download Anna and deploy your site in seconds with our quick start guide.

```text
    ___
   /   |  ____  ____  ____ _
  / /| | / __ \/ __ \/ __ `/
 / ___ |/ / / / / / / /_/ /
/_/  |_/_/ /_/_/ /_/\__,_/

A static site generator in go
```

This project was a part of the ACM PESU-ECC's yearly [AIEP](https://acmpesuecc.github.io/aiep) program, and is maintained by [Adhesh Athrey](https://github.com/DedLad), [Nathan Paul](https://github.com/polarhive), [Anirudh Sudhir](https://github.com/anirudhsudhir), and [Aditya Hegde](https://github.com/bwaklog)
