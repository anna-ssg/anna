---
date: 23-11-2023
title: XKCD Grab
draft: true
scripts:
  - dark.js
---

# xkcd-grab v2

### xkcd comics fetched using terminal 🥳

![xkcd](https://github.com/bwaklog/xkcd-grab/assets/91192289/29085083-88b4-45af-ad4c-ef6e9000150f)

This is the project i've used to demo fuzzy searching and web scraping in the HSP•PESUECC Project Expo ❤️

- [Link to Presentation 📈](https://docs.google.com/presentation/d/1KapKH3m2BfKAuwZvMP9RBlk4GlALjeME06s4oYKOkPc/edit?usp=sharing)

- [Link to Blog Post 📑](https://bwaklog.vercel.app/pasta/2023/11/18/xkcd-grab-blog-post.html)

---

Hey 👋
This is a CLI tool utilising API's for retrieving user-requested xkcd comics. Its a relatively small sized project, which is WIP cuz of a lack of data. This project is somewhat of a playground for me to explore different searching and querying techniques.

Due to data limitation, I wanted to make it a goal to make it super easy to find a specific comic based on query. The roadmap of this current is to make a smart cli tool to find the most relevant comic based on a search query.

# Table of Contents

1. [Installation and Usage](#usage)
   - [Commands Available](#commands-available)
2. [Cool Stuff](#cool-stuff)
3. [Party Feature](#party-feature)
4. [Yet To Come](#yet-to-come)
5. [Requirements - covered in Installation](#requirements)

## Installation and Usage:

<a name="usage"></a>

- Clone this repository

```bash
  git clone https://github.com/bwaklog/xkcd-grab
```

- Install requirements
  Some **pre-requisites**
  1. Python3+
  2. tesseract OCR engine that is going to be implemented in the future updates

```bash
./install.sh
```

- Add xkcd alias to the path for easier commands. Add alias to the path manually, I still have to figure out how to automate this.

```bash
alias xkcd='./xkcd.sh'
```

_Sidenote_ : the script creates a virtual env `venv`, so you might want to start using it

```bash
. venv/bin/activate
```

### Fuzzy Search Demo

![xkcd-grab-demo-fuzzy](https://github.com/bwaklog/xkcd-grab/assets/91192289/24dac00e-428b-4c06-88f0-f2ad23c3abcd)

### Web Scraping Demo

![xkcd-grab-demo-google-2](https://github.com/bwaklog/xkcd-grab/assets/91192289/d67e4689-e316-4d9b-b6d2-8907ebd3ca4d)

---

### Commands Available

<a name="commands-available"></a>

_PS_ even if u mess up the commands, there is a help file to guide you...which I am yet to complete :P

Here is a boiler plate of how the CLI commands must follow

```
xkcd <type of request> <extra commands>
```

## Cool stuff:

<a name="cool-stuff"></a>

- For MacOS systems, image is opened using the system quicklook. This has been done by utilising the `qlmanage` command
- Web Scraping uses googles best matches to find the comic you are searching for. All you need to type is a search query(anything that describes the comic)
  <br />
  <iframe src="https://i.imgur.com/xCOmCyX.mp4" allow="fullscreen" allowfullscreen="" style="height: 100%; width: 100%; aspect-ratio: 16 / 9;"></iframe>

## Yet to come

<a name="yet-to-come"></a>

This project is somewhat of a playground for me to explore different searching algorithms and querying techniques. While this might have a niche target, I want to build this tool into a more robust API client. The roadmap of this current is to make a smart cli tool to find the most relevant comic based on a search query.

The current `web scraping` function that is built into the app is the goal I am trying to achieve using data from all the 2800+ comics alone. So this is still very much a work in progress

- _Create a web interface using **flask**._..<br />
  May or may not go ahead with this option cuz the main goal was to create a cli tool. But if needed, I take a chance in making one.

  - 💾 Local Storage options for comics
  - ❤️ Creating Bookmark/Liking features
  - 📩 Creating a sharing option. Send your favorite comics to your friends with a few clicks!
  - Umm...A neat interface cuz I don't want get myself using tkinter or some other boring looking tool.

- There was supposed to be an `install.sh` script to add the `xkcd.sh` script to your alias but that didn't seem to work cuz idk how to do that
  <br /><br />

### Tabulated stuff for professionalism 🫡

| ツ  | **Feature**                                                         | **Progress**            |
| --- | ------------------------------------------------------------------- | ----------------------- |
| 🔥  | Smart Comic Search                                                  | 🕺 In progress          |
| 💾  | Local Storage Option                                                | 👍 workaround available |
| ❤️  | Liking\Bookmarking option to save comic no and not on local storage | 🔘                      |
| 📩  | Sharing feature (undecided)                                         | WAP                     |
| 🤔  | Flask generated page                                                | TBA                     |

## Party Feature:

<a name="party-feature"></a>

> ⚠️ This is very much in devlopment, but here is how you can use the little orca-mini LLM to make the cli expaliln the comic

1. Install [_ollama_](https://github.com/jmorganca/ollama)
2. Install orca-mini's LLM using ollama (about 2.0 GB)

```shell
ollama pull orca-mini
# if ur familiar with docker, you know whats going on
# also macos and linux only for now i guess (26th Nov)
```

3. Start the server in another temrinal window

```shell
ollama serve
```

4. Use the flag `-e` or `--explain` after fuzzy search, or web scraping for it to start generating after getting the results

```shell
xkcd -f -e
```

![](https://i.imgur.com/a2SsJtM.jpg)

## Requirements

<a name="requirements"></a>

What i'm using for this program:

- This isn't really a disclaimer but if you don't have quick-look (MacOS only), that's no problem! But for now all you get is:
  - 🔗 A link to the image of the post. You can open it in your default browser
  - A _very very very_ descriptive info of the post you requested for ツ
- Yeah, I haven't used this on a windows pc so far, and some of these..._**most of these**_ commands are _UNIX_ commands so, **join the the Force with BASH 🕺**
- running this in a venv for development, so do make sure you install all the requirements from `requirements.txt`
- That's it for now...nothing else to force you to install..**other than python3**🐍
- ollama and orca-mini for party feature
