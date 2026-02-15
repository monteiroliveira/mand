package pkg

import (
	"errors"
	"fmt"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg/parsers/manga"
)

type MangaParsers string

const (
	MangaDex MangaParsers = "MangaDex"
)

var SupportedMangaParsers map[string]MangaParsers = map[string]MangaParsers{
	manga.MangaDexValidLink: MangaDex,
}

func NewMangaParser(args *internal.ConsoleArgs) (manga.MangaParser, error) {
	value, ok := SupportedMangaParsers[args.Manga.Download.Source.URL.Host]
	if !ok {
		return nil, errors.Join(internal.SetSemanticError(), fmt.Errorf("Unsupported Source Link"))
	}

	switch value {
	case MangaDex:
		return manga.NewMangaDexParser(args.Manga.Download.Source.URL), nil
	default:
		return nil, errors.Join(internal.SetSemanticError(), fmt.Errorf("Cannot found parser for %s", value))
	}
}
