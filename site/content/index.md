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
  "https://adhesh.netlify.app",
  "https://sameermanvi.me",
  "https://adihegde.com",
  "https://hsp-ec.xyz",
  "https://anirudhsudhir.com",
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

// wrap the gallery in a scrollable marquee container if not already wrapped
if (!gallery.parentElement.classList.contains('marquee')) {
  const marquee = document.createElement('div');
  marquee.className = 'marquee';
  gallery.parentNode.replaceChild(marquee, gallery);
  marquee.appendChild(gallery);
  gallery.classList.add('marquee-track');
}

// setup auto-scroll + drag timeline
function setupMarquee() {
  const container = gallery.parentElement; // .marquee
  const track = gallery; // .marquee-track
  const initial = Array.from(track.children);
  if (!initial.length) return;

  // duplicate items until track width >= 2x container width
  let totalWidth = initial.reduce((sum, el) => sum + el.getBoundingClientRect().width + parseFloat(getComputedStyle(track).gap || 0), 0);
  let i = 0;
  while (totalWidth < container.getBoundingClientRect().width * 2) {
    const clone = initial[i % initial.length].cloneNode(true);
    track.appendChild(clone);
    totalWidth += clone.getBoundingClientRect().width + parseFloat(getComputedStyle(track).gap || 0);
    i++;
    if (i > 60) break;
  }

  const trackWidth = Array.from(track.children).reduce((w, node) => w + node.getBoundingClientRect().width + parseFloat(getComputedStyle(track).gap || 0), 0);
  const half = trackWidth / 2;

  // auto-scroll state
  let last = null;
  const speed = 12; // px/s
  let isPointerDown = false;
  let isHover = false;
  let isFocused = false;

  let startX = 0;
  let startScroll = 0;

  container.addEventListener('pointerdown', (e) => {
    isPointerDown = true;
    container.setPointerCapture?.(e.pointerId);
    startX = e.clientX;
    startScroll = container.scrollLeft;
  });

  container.addEventListener('pointermove', (e) => {
    if (!isPointerDown) return;
    const dx = e.clientX - startX;
    container.scrollLeft = startScroll - dx;
    if (container.scrollLeft >= half) container.scrollLeft -= half;
    if (container.scrollLeft < 0) container.scrollLeft += half;
  });

  container.addEventListener('pointerup', (e) => {
    isPointerDown = false;
    try { container.releasePointerCapture?.(e.pointerId); } catch (err) {}
  });

  container.addEventListener('pointercancel', () => { isPointerDown = false; });
  container.addEventListener('mouseleave', () => { isPointerDown = false; });

  container.addEventListener('mouseenter', () => { isHover = true; });
  container.addEventListener('mouseleave', () => { isHover = false; });
  container.addEventListener('focusin', () => { isFocused = true; });
  container.addEventListener('focusout', () => { isFocused = false; });

  function step(t) {
    if (!last) last = t;
    const dt = (t - last) / 1000;
    last = t;

    if (!isPointerDown && !isHover && !isFocused) {
      container.scrollLeft += speed * dt;
      if (container.scrollLeft >= half) container.scrollLeft -= half;
    }

    requestAnimationFrame(step);
  }

  requestAnimationFrame(step);

  // ensure initial scroll is positioned within first half
  container.scrollLeft = 0;
}

window.addEventListener('load', setupMarquee);
window.addEventListener('resize', () => {
  clearTimeout(window._marqueeTimer);
  window._marqueeTimer = setTimeout(setupMarquee, 150);
});
</script>

---

# [Quick Start Guide](./quick-start.html)

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
