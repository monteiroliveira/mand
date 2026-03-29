package scraper

import (
	"fmt"
	"strings"
)

const specialCharacters string = `.+?^$()[]{}/\|`

type RegexParser struct{}

func NewRegexParser() *RegexParser {
	return &RegexParser{}
}

func (r *RegexParser) Normalize(target string) (string, error) {
	normalized := ""
	for _, ch := range target {
		if strings.ContainsRune(specialCharacters, ch) {
			normalized += `\` + string(ch)
		} else {
			normalized += string(ch)
		}
	}
	if normalized == "" {
		return "", fmt.Errorf("Failed to normalize target")
	}
	return normalized, nil
}
