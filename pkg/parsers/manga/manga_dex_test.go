package manga

import (
	"net/url"
	"testing"

	"github.com/monteiroliveira/mand/internal/scraper/models"
)

func TestGetChapterId_ValidURL(t *testing.T) {
	u, _ := url.Parse("https://mangadex.org/chapter/abc-123-def")
	id, err := getChapterId(u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "abc-123-def" {
		t.Errorf("expected 'abc-123-def', got %q", id)
	}
}

func TestGetChapterId_NestedPath(t *testing.T) {
	u, _ := url.Parse("https://mangadex.org/chapter/some-uuid/extra")
	id, err := getChapterId(u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "extra" {
		t.Errorf("expected 'extra', got %q", id)
	}
}

func TestGetChapterId_RootPath(t *testing.T) {
	u, _ := url.Parse("https://mangadex.org/")
	id, err := getChapterId(u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// path is "/" → split gives ["", ""] → last element is ""
	if id != "" {
		t.Errorf("expected empty string for root path, got %q", id)
	}
}

func TestBuildDownloadList(t *testing.T) {
	info := &models.MangaDexChapterDownloadInfo{
		Result:  "ok",
		BaseUrl: "https://uploads.mangadex.org",
		Chapter: struct {
			Hash string   `json:"hash"`
			Data []string `json:"data"`
		}{
			Hash: "abc123",
			Data: []string{"page1.png", "page2.png", "page3.png"},
		},
	}

	links := buildDownloadList(info)
	if len(links) != 3 {
		t.Fatalf("expected 3 links, got %d", len(links))
	}

	expected := "https://uploads.mangadex.org/data/abc123/page1.png"
	if links[0] != expected {
		t.Errorf("expected %q, got %q", expected, links[0])
	}
}

func TestBuildDownloadList_Empty(t *testing.T) {
	info := &models.MangaDexChapterDownloadInfo{
		BaseUrl: "https://uploads.mangadex.org",
		Chapter: struct {
			Hash string   `json:"hash"`
			Data []string `json:"data"`
		}{
			Hash: "abc",
			Data: []string{},
		},
	}

	links := buildDownloadList(info)
	if len(links) != 0 {
		t.Errorf("expected 0 links, got %d", len(links))
	}
}

func TestNewMangaDexParser(t *testing.T) {
	u, _ := url.Parse("https://mangadex.org/chapter/test-id")
	args := &MangaParserArgs{
		Source:    u,
		Operation: DownloadOperation,
	}

	parser := NewMangaDexParser(args)
	if parser == nil {
		t.Fatal("expected non-nil parser")
	}
	if parser.source.String() != u.String() {
		t.Errorf("expected source %q, got %q", u.String(), parser.source.String())
	}
}

func TestMangaDexParser_ExtractChapterContentByList_NotSupported(t *testing.T) {
	u, _ := url.Parse("https://mangadex.org/chapter/test-id")
	args := &MangaParserArgs{
		Source:    u,
		Operation: DownloadListOperation,
	}
	parser := NewMangaDexParser(args)

	err := parser.ExtractChapterContentByList(5)
	if err == nil {
		t.Fatal("expected error for unsupported operation")
	}
	if err.Error() != "Not supported" {
		t.Errorf("expected 'Not supported', got %q", err.Error())
	}
}
