package scraper

import (
	"testing"
)

func TestRegexParserNormalize_SpecialChars(t *testing.T) {
	r := NewRegexParser()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "dot",
			input:    "mangadex.org",
			expected: `mangadex\.org`,
		},
		{
			name:     "url with dots and slashes",
			input:    "https://www.mangaread.org/manga/",
			expected: `https:\/\/www\.mangaread\.org\/manga\/`,
		},
		{
			name:     "no special chars",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "multiple special chars",
			input:    "a.b+c?d",
			expected: `a\.b\+c\?d`,
		},
		{
			name:     "parentheses and brackets",
			input:    "(test)[value]{key}",
			expected: `\(test\)\[value\]\{key\}`,
		},
		{
			name:     "pipe and caret",
			input:    "a|b^c$d",
			expected: `a\|b\^c\$d`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.Normalize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestRegexParserNormalize_EmptyString(t *testing.T) {
	r := NewRegexParser()
	_, err := r.Normalize("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}
