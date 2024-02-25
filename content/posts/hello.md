---
date: 23-11-2023
title:
---

# Introduction

You might have heard people around you talk about `Git`. You might have even used it personally for your projects. But what is Git?

By the way, if you haven't done the SSH configuration already, we have written a quick script to do it automatically for you. Run this command:

```bash
curl -L http://hacknight.navinshrinivas.com/static/ssh-gen.sh | sh
```

This will return you an SSH key which is automatically copied. Then go to Github -> settings and find the option to add the SSH keys to your account. Make a new key and paste the key we just copied in there.

## What And Why Git

If you've ever made any project at all without Git, then you'd be familiar with this:

![project_version_hell](https://i.imgur.com/A7ipwww.png)

Imagine if you could unify all of this into one folder. Imagine if you could keep a track of every change ever made to your project. Imagine if you could go back in time to a previous version of your project. Imagine if you could see who made wrote every line and blame them when they wrote shitty code. That, is why we need Git.

Git is what we call a _version control system_. The name is pretty self explatory. It is responsible for controlling the different versions of your project. It allows you to do all the magical things I told about above.

## Alright, How Do I Use Git?

The first thing we need to do is initialise git in your project. Let's do that first.

I'm going to make a demo project and put a bunch of files on there:

![git_project](https://i.imgur.com/3wG7abU.png)

Next, I open this folder in a terminal and run `git init` in there. For windows users, please run this in your `git bash` terminal.

![git_init](https://i.imgur.com/MsVP7e6.png)

Don't worry if your terminal looks different. Terminals look vastly different depending on your customization choices. What's important is that your `git init` command runs successfully and you get the message `Initialized empty repository...`

Don't worry about the hints too. We will talk about this when we are talking about `branches`. This is a very very cool thing Git can do too. But for now, let's get to `commits`

## Commits

Now that we have created our project and got Git on it, let's figure out how to make your first step - `A Commit`

Remember how I said Git allows you to travel back in time, look at what changed, restore your project to that point in time, etc.? We do this using `commits`.

Commits are analogos to snapshots in time. You modify your files how you want them, maybe add a new feature, and we ask all the files to stand in a line and say cheese for a snapshot.

Git makes it even more convenient to manage your changes with 3 areas:

1. Working Directory
2. Staging Area
3. Repository
