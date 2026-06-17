// Package renderer turns parsed SiteData into a complete static site on disk:
// HTML pages from templates, a copy of the static assets, and an RSS feed.
package renderer

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kkjorsvik/kylekjorsvik.com/internal/model"
)

const (
	// SiteTitle and friends are the site-wide constants used in templates,
	// meta tags, and the RSS feed.
	SiteTitle       = "Kyle Kjorsvik"
	SiteURL         = "https://kylekjorsvik.com"
	SiteTagline     = "Senior DevOps Engineer · Platform Engineering · Builder"
	SiteDescription = "DevOps, platform engineering, and projects."

	maxRecentPosts     = 9
	maxFeaturedProject = 6
)

// Renderer holds the configured paths and parsed templates.
type Renderer struct {
	templatesDir string
	staticDir    string
	outputDir    string
	site         model.SiteData
	buildTime    time.Time
}

// New constructs a Renderer for the given directories and parsed content.
// buildTime is stamped into the footer of every page.
func New(templatesDir, staticDir, outputDir string, site model.SiteData, buildTime time.Time) *Renderer {
	return &Renderer{
		templatesDir: templatesDir,
		staticDir:    staticDir,
		outputDir:    outputDir,
		site:         site,
		buildTime:    buildTime,
	}
}

// ctx is the data handed to every template execution. Only the fields relevant
// to a given page are populated.
type ctx struct {
	Site        model.SiteData
	Title       string
	Description string
	Nav         string // active nav key: home, blog, projects, uses, about
	URL         string // canonical path of the current page

	Post     model.Post
	Project  model.Project
	Page     model.Page
	Posts    []model.Post
	Projects []model.Project

	SiteTitle   string
	SiteURL     string
	SiteTagline string
	BuildTime   time.Time
}

func (r *Renderer) baseCtx() ctx {
	return ctx{
		Site:        r.site,
		SiteTitle:   SiteTitle,
		SiteURL:     SiteURL,
		SiteTagline: SiteTagline,
		Description: SiteDescription,
		BuildTime:   r.buildTime,
	}
}

// funcMap exposes small helpers to templates.
var funcMap = template.FuncMap{
	"year": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("2006")
	},
	"dateLong": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("January 2, 2006")
	},
	"datetime": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.UTC().Format("Jan 2, 2006 15:04 MST")
	},
	"join": strings.Join,
}

// load parses base.html together with a single page template. The page template
// redefines the "content" (and optionally "title") block declared in base.html.
func (r *Renderer) load(page string) (*template.Template, error) {
	return template.New("base.html").Funcs(funcMap).ParseFiles(
		filepath.Join(r.templatesDir, "base.html"),
		filepath.Join(r.templatesDir, "partials.html"),
		filepath.Join(r.templatesDir, page),
	)
}

// renderTo executes the given page template into outputDir/relPath, creating
// parent directories as needed.
func (r *Renderer) renderTo(page, relPath string, data ctx) error {
	tmpl, err := r.load(page)
	if err != nil {
		return fmt.Errorf("loading %s: %w", page, err)
	}

	dest := filepath.Join(r.outputDir, relPath)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "base.html", data); err != nil {
		return fmt.Errorf("rendering %s: %w", relPath, err)
	}
	return nil
}

// RenderAll renders every page of the site.
func (r *Renderer) RenderAll() error {
	if err := r.renderHome(); err != nil {
		return err
	}
	if err := r.renderPosts(); err != nil {
		return err
	}
	if err := r.renderProjects(); err != nil {
		return err
	}
	if err := r.renderPages(); err != nil {
		return err
	}
	return nil
}

func (r *Renderer) renderHome() error {
	data := r.baseCtx()
	data.Nav = "home"
	data.Title = SiteTitle
	data.URL = "/"

	data.Posts = r.site.Posts
	if len(data.Posts) > maxRecentPosts {
		data.Posts = data.Posts[:maxRecentPosts]
	}
	var featured []model.Project
	for _, p := range r.site.Projects {
		if p.Featured {
			featured = append(featured, p)
		}
	}
	if len(featured) > maxFeaturedProject {
		featured = featured[:maxFeaturedProject]
	}
	data.Projects = featured
	return r.renderTo("index.html", "index.html", data)
}

func (r *Renderer) renderPosts() error {
	// Blog index.
	data := r.baseCtx()
	data.Nav = "blog"
	data.Title = "Blog"
	data.Description = "Writing on DevOps, platform engineering, and building things."
	data.URL = "/blog/"
	data.Posts = r.site.Posts
	if err := r.renderTo("post-list.html", filepath.Join("blog", "index.html"), data); err != nil {
		return err
	}

	// Individual posts.
	for _, p := range r.site.Posts {
		d := r.baseCtx()
		d.Nav = "blog"
		d.Title = p.Title
		d.Description = p.Description
		d.URL = p.URL
		d.Post = p
		if err := r.renderTo("post.html", filepath.Join("blog", p.Slug, "index.html"), d); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) renderProjects() error {
	// Projects index.
	data := r.baseCtx()
	data.Nav = "projects"
	data.Title = "Projects"
	data.Description = "Things I'm building and have built."
	data.URL = "/projects/"
	data.Projects = r.site.Projects
	if err := r.renderTo("project-list.html", filepath.Join("projects", "index.html"), data); err != nil {
		return err
	}

	// Individual projects.
	for _, p := range r.site.Projects {
		d := r.baseCtx()
		d.Nav = "projects"
		d.Title = p.Title
		d.Description = p.Description
		d.URL = p.URL
		d.Project = p
		if err := r.renderTo("project.html", filepath.Join("projects", p.Slug, "index.html"), d); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) renderPages() error {
	for _, p := range r.site.Pages {
		d := r.baseCtx()
		d.Nav = p.Slug // "about" / "uses" line up with nav keys
		d.Title = p.Title
		d.URL = p.URL
		d.Page = p
		if err := r.renderTo("page.html", filepath.Join(p.Slug, "index.html"), d); err != nil {
			return err
		}
	}
	return nil
}

// CopyStatic copies the static asset directory into output/static.
func (r *Renderer) CopyStatic() error {
	dest := filepath.Join(r.outputDir, "static")
	return copyDir(r.staticDir, dest)
}

// WriteFile writes arbitrary bytes to a path under the output directory. Used
// by the build to drop in generated assets like the Chroma stylesheet.
func (r *Renderer) WriteFile(relPath string, content []byte) error {
	dest := filepath.Join(r.outputDir, relPath)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, content, 0o644)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// RenderFeed writes an RSS 2.0 feed at output/feed.xml from the posts.
func (r *Renderer) RenderFeed() error {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<rss version="2.0">` + "\n")
	b.WriteString("  <channel>\n")
	b.WriteString("    <title>" + escapeXML(SiteTitle) + "</title>\n")
	b.WriteString("    <link>" + escapeXML(SiteURL) + "</link>\n")
	b.WriteString("    <description>" + escapeXML(SiteDescription) + "</description>\n")

	for _, p := range r.site.Posts {
		link := SiteURL + p.URL
		b.WriteString("    <item>\n")
		b.WriteString("      <title>" + escapeXML(p.Title) + "</title>\n")
		b.WriteString("      <link>" + escapeXML(link) + "</link>\n")
		b.WriteString("      <guid>" + escapeXML(link) + "</guid>\n")
		b.WriteString("      <description>" + escapeXML(p.Description) + "</description>\n")
		b.WriteString("      <pubDate>" + p.Date.Format(time.RFC1123Z) + "</pubDate>\n")
		b.WriteString("    </item>\n")
	}

	b.WriteString("  </channel>\n")
	b.WriteString("</rss>\n")

	return r.WriteFile("feed.xml", []byte(b.String()))
}

func escapeXML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}
