package pkg

import (
	"net/url"
	"testing"

	"github.com/monteiroliveira/mand/pkg/parsers/manga"
)

func TestSupportedMangaParsers(t *testing.T) {
	if _, ok := SupportedMangaParsers["mangadex.org"]; !ok {
		t.Error("expected mangadex.org in supported parsers")
	}
	if _, ok := SupportedMangaParsers["www.mangaread.org"]; !ok {
		t.Error("expected www.mangaread.org in supported parsers")
	}
}

func TestNewMangaParser_MangaDex(t *testing.T) {
	u, _ := url.Parse("https://mangadex.org/chapter/abc123")
	args := &manga.MangaParserArgs{
		Source:    u,
		Operation: manga.DownloadOperation,
	}

	parser, err := NewMangaParser(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parser == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestNewMangaParser_MangaRead(t *testing.T) {
	u, _ := url.Parse("https://www.mangaread.org/manga/one-piece/chapter-1/")
	args := &manga.MangaParserArgs{
		Source:    u,
		Operation: manga.DownloadOperation,
	}

	parser, err := NewMangaParser(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parser == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestNewMangaParser_UnsupportedSource(t *testing.T) {
	u, _ := url.Parse("https://unknown-manga-site.com/chapter/1")
	args := &manga.MangaParserArgs{
		Source:    u,
		Operation: manga.DownloadOperation,
	}

	_, err := NewMangaParser(args)
	if err == nil {
		t.Fatal("expected error for unsupported source")
	}
}

func TestNewMangaParser_InvalidBatchSize(t *testing.T) {
	u, _ := url.Parse("https://www.mangaread.org/manga/test/")
	args := &manga.MangaParserArgs{
		Source:        u,
		Operation:     manga.DownloadListOperation,
		ListBatchSize: 0,
	}

	_, err := NewMangaParser(args)
	if err == nil {
		t.Fatal("expected error for invalid batch size")
	}
}

func TestExecute_InvalidOperation(t *testing.T) {
	u, _ := url.Parse("https://mangadex.org/chapter/abc123")
	args := &manga.MangaParserArgs{
		Source:    u,
		Operation: manga.Operation(99),
	}

	parser, err := NewMangaParser(args)
	if err != nil {
		t.Fatalf("unexpected error creating parser: %v", err)
	}

	err = Execute(parser, args)
	if err == nil {
		t.Fatal("expected error for invalid operation")
	}
}
