# Project Rebuild 29th March

## New proposed project structure

- Modules previously part of `cmd/anna/utils.go` and `cmd/anna/main.go` are to be split between `pkg/parsers/`, `pkg/engine/` and `pkg/helper`

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

## New proposed rendering system

- Currently there are two separate types of files that are rendered. One set are user defined files whicha are for example docs.md, timeline.md... These are specific to a user.

- The second set of files that are rendered are files that belong to `tags.html`, `sub-tags.html` and `posts.html`

- Splitting the rendering system to make parallelisation easier. Now the generator/engine has a method to render "anna specific" pages and another method to rendere "user defined" pages which include all the user pages and posts

![new proposed project rendering structure](https://i.imgur.com/LgCDh4P.png)

## Other restructuring

- [x] `main.go` will only handle flags
- [x] Split Generator struct
- [x] Split Rendering system
- [ ] Improve SEO

## Tests

- Write tests for various components of the ss
  - [ ] tests for `helper` package
  - [x] tests for `engine` package
  - [x] tests for `parser` package
  - [ ] tests for cmd/ package
  - [ ] tests for `main.go`
- [x] Add status checks for PRs and fix netlify deploy preview
- [ ] Minimum test coverage of 60% for all packages
  - [ ] Minimum test coverage of 60% for `engine`
  - [ ] Minimum test coverage of 60% for `helper`
  - [x] Minimum test coverage of 60% for `parser`
