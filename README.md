# kylekjorsvik.com

A custom static site generator in Go for my personal portfolio and blog. No
framework, no client-side rendering — Markdown in, static HTML + CSS out.

## How it works

`cmd/build` walks `content/`, parses Markdown with YAML frontmatter, renders the
`templates/` into `output/`, copies `static/`, generates a Chroma syntax stylesheet
for both themes, and writes an RSS 2.0 feed.

- **Parsing** — `internal/parser` (goldmark + goldmark-highlighting, yaml.v3)
- **Rendering** — `internal/renderer` (html/template)
- **Models** — `internal/model`

## Build

```sh
go build -o build ./cmd/build/
./build
```

Output lands in `output/` (gitignored). Serve it with anything static:

```sh
cd output && python3 -m http.server 8080
```

## Content

| Type     | Location             | URL                   |
|----------|----------------------|-----------------------|
| Posts    | `content/posts/`     | `/blog/<slug>/`       |
| Projects | `content/projects/`  | `/projects/<slug>/`   |
| Pages    | `content/pages/`     | `/<slug>/`            |

Slugs come from the filename. Set `draft: true` in frontmatter to exclude a post
or project from the build. See existing files for the frontmatter fields.

## Theming

Dark-mode-first with a light toggle. Colors live in CSS variables in
`static/css/custom.css`, swapped via `data-theme` on `<html>` and persisted to
`localStorage`. Code blocks use the `dracula` (dark) and `github` (light) Chroma
styles, generated into `static/css/syntax.css` at build time and scoped by the same
`data-theme` attribute.

## Deploy

`.woodpecker.yml` builds on push to `main` and rsyncs `output/` to the Caddy-served
directory on the remote host. Deploy secrets: `DEPLOY_HOST`, `DEPLOY_USER`,
`DEPLOY_KEY`, `DEPLOY_PATH`.
