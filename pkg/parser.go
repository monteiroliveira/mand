package pkg

import (
	"errors"
	"fmt"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg/parsers"
)

type Parsers string

const (
	MangaDex Parsers = "MangaDex"
)

var SupportedParsers map[string]Parsers = map[string]Parsers{
	parsers.MangaDexValidLink: MangaDex,
}

func NewParser(args *internal.ConsoleArgs) (parsers.Parser, error) {
	value, ok := SupportedParsers[args.Source.URL.Host]
	if !ok {
		return nil, errors.Join(internal.SetSemanticError(), fmt.Errorf("Unsupported Source Link"))
	}

	switch value {
	case MangaDex:
		return parsers.NewMangaDexParser(args.Source.URL), nil
	default:
		return nil, errors.Join(internal.SetSemanticError(), fmt.Errorf("Cannot found parser for %s", value))
	}
}
