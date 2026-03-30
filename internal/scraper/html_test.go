package scraper

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func parseHTML(t *testing.T, s string) *html.Node {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	return doc
}

func TestFindHtmlContent_MetaOgTitle(t *testing.T) {
	h := NewHtmlManager()
	doc := parseHTML(t, `<html><head>
		<meta property="og:title" content="Chapter 42 - One Piece"/>
	</head><body></body></html>`)

	result := h.FindHtmlContent(doc, "meta", "property", "^og:title$")
	if result != "Chapter 42 - One Piece" {
		t.Errorf("expected 'Chapter 42 - One Piece', got %q", result)
	}
}

func TestFindHtmlContent_NotFound(t *testing.T) {
	h := NewHtmlManager()
	doc := parseHTML(t, `<html><head></head><body><p>hello</p></body></html>`)

	result := h.FindHtmlContent(doc, "meta", "property", "og:title")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestFindHtmlContent_ImgSrc(t *testing.T) {
	h := NewHtmlManager()
	doc := parseHTML(t, `<html><body>
		<img id="image-0" src="https://example.com/page1.jpg"/>
	</body></html>`)

	result := h.FindHtmlContent(doc, "img", "id", "^image-0$")
	if result != "https://example.com/page1.jpg" {
		t.Errorf("expected image URL, got %q", result)
	}
}

func TestFindHtmlContentData_H1WithId(t *testing.T) {
	h := NewHtmlManager()
	doc := parseHTML(t, `<html><body>
		<h1 id="chapter-heading">Chapter 1 - Romance Dawn</h1>
	</body></html>`)

	result := h.FindHtmlContentData(doc, "h1", "id", "chapter-heading")
	if result != "Chapter 1 - Romance Dawn" {
		t.Errorf("expected 'Chapter 1 - Romance Dawn', got %q", result)
	}
}

func TestFindHtmlContentData_NotFound(t *testing.T) {
	h := NewHtmlManager()
	doc := parseHTML(t, `<html><body>
		<h1 id="other-heading">Something</h1>
	</body></html>`)

	result := h.FindHtmlContentData(doc, "h1", "id", "chapter-heading")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestListHtmlContent_MultipleLinks(t *testing.T) {
	h := NewHtmlManager()
	doc := parseHTML(t, `<html><body>
		<a href="https://example.com/manga/chapter-1">Ch 1</a>
		<a href="https://example.com/manga/chapter-2">Ch 2</a>
		<a href="https://example.com/other">Other</a>
	</body></html>`)

	results := h.ListHtmlContent(doc, "a", "href", "chapter-.*")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d: %v", len(results), results)
	}
}

func TestListHtmlContent_NoMatches(t *testing.T) {
	h := NewHtmlManager()
	doc := parseHTML(t, `<html><body>
		<a href="https://example.com/about">About</a>
	</body></html>`)

	results := h.ListHtmlContent(doc, "a", "href", "chapter-.*")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFindHtmlContent_NilNode(t *testing.T) {
	h := NewHtmlManager()
	result := h.FindHtmlContent(nil, "meta", "property", "og:title")
	if result != "" {
		t.Errorf("expected empty string for nil node, got %q", result)
	}
}

func TestListHtmlContent_NilNode(t *testing.T) {
	h := NewHtmlManager()
	results := h.ListHtmlContent(nil, "a", "href", "chapter-.*")
	if results != nil {
		t.Errorf("expected nil for nil node, got %v", results)
	}
}
