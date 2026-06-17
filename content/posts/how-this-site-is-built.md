---
title: "How This Site Is Built"
tags: ["go", "devops", "meta"]
draft: false
description: "A custom static site generator in Go, deployed via Woodpecker CI on merge to main. No frameworks, no CMS, just goldmark and html/template."
---

Most developers building a personal site reach for Hugo, Astro, or a hosted platform.
I built my own static site generator in Go. Not because the alternatives are bad — Hugo
in particular is excellent — but because writing your own is a more interesting problem,
and the result is something I understand completely.

## Why custom

I've used Hugo before. It's fast, well-documented, and has a large theme ecosystem. But
for a personal site I'm going to maintain indefinitely, I'd rather own a few hundred
lines of Go I wrote than debug someone else's template system when something breaks.

There's also the honest answer: I'm actively learning Go, and building something real
is how I learn. A static site generator is a well-scoped project with clear requirements,
real output you can look at, and enough interesting problems (frontmatter parsing,
template composition, RSS generation) to make it worth doing.

## How it works

The generator lives in `cmd/build/` and does five things:

1. Walk `content/posts/`, `content/projects/`, and `content/pages/` and parse each
   Markdown file
2. Extract YAML frontmatter (title, tags, description, draft status) using
   `gopkg.in/yaml.v3`
3. Render Markdown to HTML using [goldmark](https://github.com/yuin/goldmark) with
   syntax highlighting via goldmark-highlighting and Chroma
4. Apply `html/template` layouts and write static HTML to `output/`
5. Generate an RSS 2.0 feed at `/feed.xml`

That's it. No incremental builds, no asset pipeline, no plugin system. The whole thing
builds in under a second.

## Content

Every post, project page, and static page is a Markdown file with YAML frontmatter:

```yaml
---
title: "How This Site Is Built"
tags: ["go", "devops", "meta"]
draft: false
description: "A short description for the post list and meta tags."
---
```

Setting `draft: true` skips the file during the build. That's the entire content
management system.

## Deployment

The site is hosted on an EC2 instance served by nginx. Deployment is handled by
[Woodpecker CI](https://woodpecker-ci.org/) — my self-hosted CI runner — on every
push to main. The pipeline builds the Go binary, runs the build, and rsync's the
`output/` directory to the server. The whole deploy takes about 15 seconds.

Forgejo is the source of truth for the repository, mirrored to GitHub as a backup.
A push to main is all it takes to ship a change — including new posts — so deploying
is just part of my normal git workflow.

## Design

The site uses [Pico CSS](https://picocss.com/) as a baseline with custom CSS variables
on top. Dark mode by default, light mode toggle persisted to localStorage. The color
palette is zinc and cyan — dark backgrounds, cards that sit slightly above the page
surface, and a cyan accent that shows up on links and active nav items.

No JavaScript frameworks. The only script on the page is the theme toggle, which is
about ten lines inlined in the `<head>` to prevent a flash of the wrong theme on load.

## What's missing

No pagination yet — I'll add it when the post count warrants it. No search, no
comments, no analytics. The RSS feed is there if you want to follow along without
giving anything up.
