// Package parser reads Markdown content files with YAML frontmatter and
// turns them into the typed structs in internal/model. It owns the configured
// goldmark instance (with Chroma syntax highlighting) so rendering is
// consistent across every content type.
package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"gopkg.in/yaml.v3"

	"github.com/kkjorsvik/kylekjorsvik.com/internal/model"
)

// md is the shared goldmark instance. Syntax highlighting emits CSS classes
// (not inline styles) so we can ship one stylesheet per theme and swap them
// with the same data-theme attribute the rest of the site uses.
var md = goldmark.New(
	goldmark.WithExtensions(
		highlighting.NewHighlighting(
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
			),
		),
	),
)

// frontmatter is the union of every field we read from any content type. Each
// content type only fills the fields relevant to it.
type frontmatter struct {
	Title       string       `yaml:"title"`
	Tags        []string     `yaml:"tags"`
	Draft       bool         `yaml:"draft"`
	Description string       `yaml:"description"`
	Status      string       `yaml:"status"`
	Date        string       `yaml:"date"`
	Featured    bool         `yaml:"featured"`
	Links       []model.Link `yaml:"links"`
}

// parseFile reads a content file, splits the YAML frontmatter from the body,
// unmarshals the frontmatter, and renders the body to HTML. It also returns
// the file modification time, used as a fallback publish date.
func parseFile(path string) (frontmatter, template.HTML, time.Time, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return frontmatter{}, "", time.Time{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return frontmatter{}, "", time.Time{}, err
	}

	fmBytes, bodyBytes, err := splitFrontmatter(raw)
	if err != nil {
		return frontmatter{}, "", time.Time{}, fmt.Errorf("%s: %w", path, err)
	}

	var fm frontmatter
	if len(fmBytes) > 0 {
		if err := yaml.Unmarshal(fmBytes, &fm); err != nil {
			return frontmatter{}, "", time.Time{}, fmt.Errorf("%s: parsing frontmatter: %w", path, err)
		}
	}

	var buf bytes.Buffer
	if err := md.Convert(bodyBytes, &buf); err != nil {
		return frontmatter{}, "", time.Time{}, fmt.Errorf("%s: rendering markdown: %w", path, err)
	}

	return fm, template.HTML(buf.String()), info.ModTime(), nil
}

// splitFrontmatter separates a leading `---`-delimited YAML block from the
// Markdown body. A file with no frontmatter block is returned as all body.
func splitFrontmatter(raw []byte) (fm, body []byte, err error) {
	text := string(raw)
	// Normalise to make delimiter matching simpler.
	if !strings.HasPrefix(text, "---\n") && !strings.HasPrefix(text, "---\r\n") {
		return nil, raw, nil
	}

	// Drop the opening delimiter line.
	rest := strings.TrimPrefix(text, "---")
	rest = strings.TrimPrefix(rest, "\r")
	rest = strings.TrimPrefix(rest, "\n")

	// Find the closing delimiter at the start of a line.
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return nil, nil, fmt.Errorf("frontmatter opened but never closed")
	}

	fm = []byte(rest[:idx])
	after := rest[idx+len("\n---"):]
	// Trim the remainder of the closing delimiter line.
	if nl := strings.IndexByte(after, '\n'); nl != -1 {
		after = after[nl+1:]
	} else {
		after = ""
	}
	return fm, []byte(after), nil
}

// slugFromPath derives a URL slug from a file path: the base name minus .md.
func slugFromPath(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// ParsePost parses a single blog post file.
func ParsePost(path string) (model.Post, error) {
	fm, body, mtime, err := parseFile(path)
	if err != nil {
		return model.Post{}, err
	}
	slug := slugFromPath(path)
	return model.Post{
		Title:       fm.Title,
		Slug:        slug,
		Tags:        fm.Tags,
		Draft:       fm.Draft,
		Description: fm.Description,
		Body:        body,
		URL:         "/blog/" + slug + "/",
		Date:        resolveDate(fm.Date, mtime),
	}, nil
}

// ParseProject parses a single project file.
func ParseProject(path string) (model.Project, error) {
	fm, body, _, err := parseFile(path)
	if err != nil {
		return model.Project{}, err
	}
	slug := slugFromPath(path)
	return model.Project{
		Title:       fm.Title,
		Slug:        slug,
		Description: fm.Description,
		Status:      fm.Status,
		Tags:        fm.Tags,
		Links:       fm.Links,
		Draft:       fm.Draft,
		Featured:    fm.Featured,
		Body:        body,
		URL:         "/projects/" + slug + "/",
	}, nil
}

// ParsePage parses a single static page file.
func ParsePage(path string) (model.Page, error) {
	fm, body, _, err := parseFile(path)
	if err != nil {
		return model.Page{}, err
	}
	slug := slugFromPath(path)
	return model.Page{
		Title: fm.Title,
		Slug:  slug,
		Body:  body,
		URL:   "/" + slug + "/",
	}, nil
}

// resolveDate uses the frontmatter date if present and parseable, otherwise
// falls back to the file modification time.
func resolveDate(raw string, mtime time.Time) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return mtime
	}
	for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t
		}
	}
	return mtime
}

// ChromaCSS returns combined syntax-highlighting CSS for both themes.
//
// Chroma v2 scopes its generated rules by the style's mode (e.g. `.chroma.dark
// .k`) and bakes that same mode class into the rendered wrapper at build time.
// That can't be toggled at runtime, so we strip the mode qualifier to get plain
// `.chroma .k` selectors. The dark style is then emitted as the default, and the
// light style is scoped under `[data-theme="light"]` so it overrides only in
// light mode — matching the theme attribute used across the rest of the site.
func ChromaCSS(darkStyle, lightStyle string) (string, error) {
	formatter := chromahtml.New(chromahtml.WithClasses(true))

	dark := styles.Get(darkStyle)
	if dark == nil {
		return "", fmt.Errorf("unknown chroma style %q", darkStyle)
	}
	light := styles.Get(lightStyle)
	if light == nil {
		return "", fmt.Errorf("unknown chroma style %q", lightStyle)
	}

	var darkBuf, lightBuf bytes.Buffer
	if err := formatter.WriteCSS(&darkBuf, dark); err != nil {
		return "", err
	}
	if err := formatter.WriteCSS(&lightBuf, light); err != nil {
		return "", err
	}

	var out strings.Builder
	out.WriteString("/* Dark syntax theme: ")
	out.WriteString(darkStyle)
	out.WriteString(" */\n")
	out.WriteString(stripMode(darkBuf.String(), dark.Mode().String()))
	out.WriteString("\n/* Light syntax theme: ")
	out.WriteString(lightStyle)
	out.WriteString(" */\n")
	out.WriteString(scopeCSS(stripMode(lightBuf.String(), light.Mode().String()), `[data-theme="light"] `))
	return out.String(), nil
}

// stripMode removes the style's mode qualifier (e.g. ".dark") from the
// `.chroma` and `.bg` selectors so the rules match any element carrying the
// `chroma` class, regardless of the mode class baked into the wrapper HTML.
func stripMode(css, mode string) string {
	if mode == "" {
		return css
	}
	css = strings.ReplaceAll(css, ".chroma."+mode, ".chroma")
	css = strings.ReplaceAll(css, ".bg."+mode, ".bg")
	return css
}

// scopeCSS prefixes every selector in a block of CSS rules with the given
// ancestor selector. It assumes simple Chroma output: one rule per line of the
// form `selector { ... }`, possibly preceded by a `/* comment */`.
func scopeCSS(css, prefix string) string {
	var out strings.Builder
	for _, line := range strings.Split(css, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			out.WriteByte('\n')
			continue
		}
		// Preserve any leading comment, then prefix the selector portion.
		rest := trimmed
		if strings.HasPrefix(rest, "/*") {
			if end := strings.Index(rest, "*/"); end != -1 {
				out.WriteString(rest[:end+2])
				out.WriteByte(' ')
				rest = strings.TrimSpace(rest[end+2:])
			}
		}
		if rest == "" {
			out.WriteByte('\n')
			continue
		}
		out.WriteString(prefix)
		out.WriteString(rest)
		out.WriteByte('\n')
	}
	return out.String()
}
