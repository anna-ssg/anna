---
title: Implementing Zettelkasten in Anna
description: This post focuses of the Proof of Concept behind how we plan to integrate zettelkasten in anna(our SSG) and implementing the new note taking functionality, supporting a new version of *Deep Data Merge* along with
authors:
  - Aditya Hegde
  - Anirudh Sudhir
date: 2024-04-10
draft: false
tags:
  - blog
  - tech
  - aiep
collections: ["posts", "anna"]
---

# Proof of concept

[Andy Matuschak's](https://notes.andymatuschak.org/) working notes is the key inspiration for this concept.
We are trying to deviate from the "general idea" of a blog site and focus more on this niche use case.
By integrating this feature, we are letting users to create a space to store there "zettels" and share their short notes as well.

We aren't trying to re-invent the process of making an editor that helps users maintain these zettels as there are already some fantastic applications,
namely [Obsidian](https://obsidian.md/), [Ginko Writer](https://app.gingkowriter.com) and [Evergreen Notes](https://evergreennotes.com/).
Our application as a rather needs to provide a generator to stitch these notes together to make it accessible on the site.

## Stages of Generations of Notes

---

## 1.0 Figuring out the named links.

All notes usually have titles as a phrase that can be referred to in a certain note. Our job as the SSG is to link these two notes together. For example

```md
Sed ut velit ante. Suspendisse ac porta urna, eget iaculis dui. Lorem ipsum
dolor sit amet, consectetur adipiscing elit. Donec vel enim dolor.
[[Nunc ullamcorper]] neque ut mattis commodo. Morbi bibendum sem accumsan mi
imperdiet, id egestas nulla posuere. Morbi sodales justo euismod nulla
porttitor lobortis sed ut sem.
```

The above md file is referencing a note namely `Nunc ullamcorper`. What needs to be done is, this "callout" is to be replaced by a link to that specific "note".

### 1.1 Basic Working Model

The parser must search through the `Body` section of these _notes_.
There are supposed to be "user defined" references to notes, which the parser must identify and add.
The specific reference to the `template data` of that specific post is appended with the information of all the links that it has found during parsing of the file.
This can be utilised later by the templating engine.

### <a name="linking-automation"></a>1.2 Automation of linking process

The previous method suggests the user has to manually link posts.
With automation, the goal is to remove the need for manually entering these links.
Instead we plan to use the `[[]]` callouts to the note name.
For example, `[[Nunc ullamcorper]]` will reference the markdown file which contains "Nunc ullamcoper" as the _Title_.
These callouts to other notes are to be picked out by the parser and replaced with a proper markdown reference in the buffer, so that the acutal file remains untouched.
For example `[[Nunc ullamcorper]]` will be updated to `[Nunc ullamcorper](/notes/zettel_name/123782734234)`.

### 1.3 Automation of file creation

This is a step forward from from [Automation of linking process](#linking-automation). As we are not a text editing application, this feature will make the process of creating subnotes simpler.

Other than just the parser identifying the `[[]]` callouts, during the live reload process, an additional file will be created under the same _zettel_, provided a name is mentioned in the _callout_.

---

## 2.0 Restructuring

As of now, our content directory looks somewhat like this:

```text
|— pages markdown files
|— /posts : post dir containing markdown files for all posts
```

To this, we plan to add an extra directory named `notes` that will handle all of our zettles.
Each zettel (related notes) can be organised in its own sub directories.

> Reference:
>
> - [Zettelkasten creation](https://zettelkasten.de/posts/create-zettel-from-reading-notes/)
> - [Evergreen Notes should be Atomic](https://notes.andymatuschak.org/Evergreen_notes_should_be_atomic)

We expect users users to specify the head of these zettels by themseves in the frontmatter explicitly

```yaml
---
title: Note taking can be fun
date: 2024-04-08
type: zettel
head: true
---
```

### 2.1 Concept of the `Mega Struct` (Deep Data Merge)

As each zettel must have access to the information of all other zettels, the implementation of a Deep Data Merge is quite necessary.
Each page is rendered by passsing a `Mega Struct` that the entire data of the notes section.
This struct will have the following fields:

```go
type NotesMerged struct {
	//Stores all the notes
	Notes map[template.URL]Note

	//Stores the links of each note to other notes
	LinkStore map[string][]*NoteStruct
}
```

`LinkStore` is a map which contains a slice of pointers to the _linked notes_ which eliminates data redundancy to certain extent.
This is an essential feature as Zettel emphasises on dense linking of notes.

This `LinkStore` Map is the second step of generation after all the notes in the `notes` directory have been successfully parsed.
Once the link maps have been generated, we use a similar render note function to produce the linked html files.

Each Note can is a struct that stores all of the data of a particular note, including the frontmatter.

```go
type Note struct {
	// Note data including frontmatter and content
	LinkedNotes		[]string
}
```

---

# TODO for zettelkasten impl

- [x] Generation of Linked Notes
  - [x] Implement 1.1 version of linking (user defined references to notes)
  - [x] Implement automation for the process of linking. Using `[[]]` callouts to file names.
- Tests:
  - [ ] unit tests for parsing parocess of the package rendering processes of package
