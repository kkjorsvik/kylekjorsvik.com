// Command build is the static site generator for kylekjorsvik.com. It parses
// Markdown content into typed models, renders templates into output/, copies
// static assets, generates the Chroma syntax-highlighting stylesheet, and
// writes the RSS feed.
package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/kkjorsvik/kylekjorsvik.com/internal/model"
	"github.com/kkjorsvik/kylekjorsvik.com/internal/parser"
	"github.com/kkjorsvik/kylekjorsvik.com/internal/renderer"
)

const (
	contentDir   = "content"
	templatesDir = "templates"
	staticDir    = "static"
	outputDir    = "output"

	darkChromaStyle  = "dracula"
	lightChromaStyle = "github"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("build failed: %v", err)
	}
}

func run() error {
	site, err := loadContent()
	if err != nil {
		return err
	}

	// Start from a clean output directory so deletes propagate.
	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	r := renderer.New(templatesDir, staticDir, outputDir, site, time.Now().UTC())

	if err := r.RenderAll(); err != nil {
		return err
	}
	if err := r.CopyStatic(); err != nil {
		return err
	}

	// Generate the syntax-highlighting stylesheet for both themes.
	css, err := parser.ChromaCSS(darkChromaStyle, lightChromaStyle)
	if err != nil {
		return err
	}
	if err := r.WriteFile(filepath.Join("static", "css", "syntax.css"), []byte(css)); err != nil {
		return err
	}

	if err := r.RenderFeed(); err != nil {
		return err
	}

	log.Printf("built %d posts, %d projects, %d pages",
		len(site.Posts), len(site.Projects), len(site.Pages))
	return nil
}

// loadContent walks the content directories and parses everything, skipping
// drafts for posts and projects.
func loadContent() (model.SiteData, error) {
	var site model.SiteData

	postFiles, err := markdownFiles(filepath.Join(contentDir, "posts"))
	if err != nil {
		return site, err
	}
	for _, f := range postFiles {
		p, err := parser.ParsePost(f)
		if err != nil {
			return site, err
		}
		if p.Draft {
			continue
		}
		site.Posts = append(site.Posts, p)
	}

	projectFiles, err := markdownFiles(filepath.Join(contentDir, "projects"))
	if err != nil {
		return site, err
	}
	for _, f := range projectFiles {
		p, err := parser.ParseProject(f)
		if err != nil {
			return site, err
		}
		if p.Draft {
			continue
		}
		site.Projects = append(site.Projects, p)
	}

	pageFiles, err := markdownFiles(filepath.Join(contentDir, "pages"))
	if err != nil {
		return site, err
	}
	for _, f := range pageFiles {
		p, err := parser.ParsePage(f)
		if err != nil {
			return site, err
		}
		site.Pages = append(site.Pages, p)
	}

	// Newest posts first.
	sort.SliceStable(site.Posts, func(i, j int) bool {
		return site.Posts[i].Date.After(site.Posts[j].Date)
	})
	// Stable, predictable ordering for projects and pages.
	sort.SliceStable(site.Projects, func(i, j int) bool {
		return site.Projects[i].Title < site.Projects[j].Title
	})
	sort.SliceStable(site.Pages, func(i, j int) bool {
		return site.Pages[i].Title < site.Pages[j].Title
	})

	return site, nil
}

// markdownFiles returns the sorted list of *.md files directly under dir. A
// missing directory yields an empty list rather than an error.
func markdownFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".md" {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)
	return files, nil
}
