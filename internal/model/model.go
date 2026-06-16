// Package model holds the shared data structures used across the site
// generator: parsed content (posts, projects, pages) and the aggregate
// SiteData passed into templates.
package model

import (
	"html/template"
	"time"
)

// Link is a labelled external URL shown on a project page.
type Link struct {
	Label string
	URL   string
}

// Post is a single blog post parsed from content/posts/*.md.
type Post struct {
	Title       string
	Slug        string
	Tags        []string
	Draft       bool
	Description string // used in meta tags and the post list
	Body        template.HTML
	URL         string    // e.g. /blog/my-post-title/
	Date        time.Time // from frontmatter, else file mtime
}

// Project is a single project page parsed from content/projects/*.md.
type Project struct {
	Title       string
	Slug        string
	Description string
	Status      string // e.g. "active", "paused", "complete"
	Tags        []string
	Links       []Link
	Draft       bool
	Featured    bool // surfaced in the homepage "Featured Projects" section
	Body        template.HTML
	URL         string
}

// Page is a generic static page parsed from content/pages/*.md.
type Page struct {
	Title string
	Slug  string
	Body  template.HTML
	URL   string
}

// SiteData is the aggregate of all parsed content, handed to templates.
type SiteData struct {
	Posts    []Post
	Projects []Project
	Pages    []Page
}
