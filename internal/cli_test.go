package internal

import (
	"testing"
)

func TestSourceUnmarshalText_ValidURL(t *testing.T) {
	s := &Source{}
	err := s.UnmarshalText([]byte("https://mangadex.org/chapter/abc123"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.URL == nil {
		t.Fatal("expected URL to be set")
	}
	if s.URL.Host != "mangadex.org" {
		t.Errorf("expected host 'mangadex.org', got %q", s.URL.Host)
	}
	if s.URL.Scheme != "https" {
		t.Errorf("expected scheme 'https', got %q", s.URL.Scheme)
	}
}

func TestSourceUnmarshalText_MissingScheme(t *testing.T) {
	s := &Source{}
	err := s.UnmarshalText([]byte("mangadex.org/chapter/abc123"))
	if err == nil {
		t.Fatal("expected error for URL without scheme")
	}
}

func TestSourceUnmarshalText_MissingHost(t *testing.T) {
	s := &Source{}
	err := s.UnmarshalText([]byte("https://"))
	if err == nil {
		t.Fatal("expected error for URL without host")
	}
}

func TestSourceUnmarshalText_EmptyString(t *testing.T) {
	s := &Source{}
	err := s.UnmarshalText([]byte(""))
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

func TestSourceUnmarshalText_FullURL(t *testing.T) {
	s := &Source{}
	err := s.UnmarshalText([]byte("https://www.mangaread.org/manga/one-piece/chapter-1/"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.URL.Host != "www.mangaread.org" {
		t.Errorf("expected host 'www.mangaread.org', got %q", s.URL.Host)
	}
	if s.URL.Path != "/manga/one-piece/chapter-1/" {
		t.Errorf("unexpected path: %q", s.URL.Path)
	}
}
