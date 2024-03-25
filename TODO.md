- New file strucure for packages

```
pkg
├─── helpers
│   ├─── helpers.go
├─── generators
│   ├─── generators.go
├─── parsers
    ├─── parsers.go

```

- modules previously part of `cmd/anna/utils.go` are to be split between `pkg/parsers/parsers.go` and `pkg/generators/generators.go`

## New proposed rendering system

- Currently there are two separate types of files that are rendered. One set are user defined files whicha are for example docs.md, timeline.md... These are specific to a user.

- The second set of files that are rendered are files that belong to `tags.html`, `sub-tags.html` and `posts.html`

![new proposed project rendering structure](https://i.imgur.com/LgCDh4P.png)

## Other restructuring

- `main.go` will only handle flags

- [ ] Split generator struct
- [ ] Fix SEO
