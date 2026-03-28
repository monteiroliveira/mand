package pkg

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg/parsers/manga"
)

type MangaParsers string

const (
	MangaDex  MangaParsers = "MangaDex"
	MangaRead MangaParsers = "MangaRead"
)

var SupportedMangaParsers map[string]MangaParsers = map[string]MangaParsers{
	manga.MangaDexValidLink:  MangaDex,
	manga.MangaReadValidLink: MangaRead,
}

func NewMangaParser(source *url.URL) (manga.MangaParser, error) {
	value, ok := SupportedMangaParsers[source.Host]
	if !ok {
		return nil, errors.Join(internal.SetSemanticError(), fmt.Errorf("Unsupported Source Link"))
	}

	switch value {
	case MangaDex:
		return manga.NewMangaDexParser(source), nil
	case MangaRead:
		return manga.NewMangaReadParser(source), nil
	default:
		return nil, errors.Join(internal.SetSemanticError(), fmt.Errorf("Cannot found parser for %s", value))
	}
}
