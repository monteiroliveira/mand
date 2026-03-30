package manga

import (
	"net/url"
	"testing"
)

func TestNewMangaReadParser(t *testing.T) {
	u, _ := url.Parse("https://www.mangaread.org/manga/one-piece/chapter-1")
	args := &MangaParserArgs{
		Source:    u,
		Operation: DownloadOperation,
	}

	parser := NewMangaReadParser(args)
	if parser == nil {
		t.Fatal("expected non-nil parser")
	}
	if parser.args.Source.String() != u.String() {
		t.Errorf("expected source %q, got %q", u.String(), parser.args.Source.String())
	}
}

func TestNewMangaReadParser_FieldsInitialized(t *testing.T) {
	u, _ := url.Parse("https://www.mangaread.org/manga/one-piece/chapter-1")
	args := &MangaParserArgs{
		Source:    u,
		Operation: DownloadOperation,
	}

	parser := NewMangaReadParser(args)
	if parser.client == nil {
		t.Error("expected non-nil client")
	}
	if parser.regex == nil {
		t.Error("expected non-nil regex")
	}
	if parser.imageManager == nil {
		t.Error("expected non-nil imageManager")
	}
	if parser.htmlManager == nil {
		t.Error("expected non-nil htmlManager")
	}
}

func TestMangaReadValidLink(t *testing.T) {
	if MangaReadValidLink != "www.mangaread.org" {
		t.Errorf("expected 'www.mangaread.org', got %q", MangaReadValidLink)
	}
}

func TestMangaReadParser_ImplementsMangaParser(t *testing.T) {
	u, _ := url.Parse("https://www.mangaread.org/manga/one-piece/chapter-1")
	args := &MangaParserArgs{
		Source:    u,
		Operation: DownloadOperation,
	}

	var _ MangaParser = NewMangaReadParser(args)
}
