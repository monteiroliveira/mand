package pkg

import (
	"errors"
	"fmt"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg/parsers/manga"
)

type MangaParsers string

const (
	MangaDex  MangaParsers = "MangaDex"
	MangaRead MangaParsers = "MangaRead"
)

// TODO: change valid link to regexp
var SupportedMangaParsers map[string]MangaParsers = map[string]MangaParsers{
	manga.MangaDexValidLink:  MangaDex,
	manga.MangaReadValidLink: MangaRead,
}

func NewMangaParser(args *manga.MangaParserArgs) (manga.MangaParser, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	value, ok := SupportedMangaParsers[args.Source.Host]
	if !ok {
		return nil, errors.Join(internal.SetSemanticError(), fmt.Errorf("Unsupported Source Link"))
	}

	switch value {
	case MangaDex:
		return manga.NewMangaDexParser(args), nil
	case MangaRead:
		return manga.NewMangaReadParser(args), nil
	default:
		return nil, errors.Join(internal.SetSemanticError(), fmt.Errorf("Cannot found parser for %s", value))
	}
}

func Execute(parser manga.MangaParser, args *manga.MangaParserArgs) error {
	switch args.Operation {
	case manga.DownloadOperation:
		pages, err := parser.ExtractChapterContent()
		if err != nil {
			return fmt.Errorf("Error in source content extraction, get %s\n", err)
		}
		chn, err := parser.ExtractChapterName()
		if err != nil {
			return fmt.Errorf("Error in get chapter name, get %s\n", err)
		}
		err = parser.DownloadPages(pages, chn)
		if err != nil {
			return fmt.Errorf("Error in pages download, get %s\n", err)
		}
	case manga.DownloadListOperation:
		err := parser.ExtractChapterContentByList(args.ListBatchSize)
		if err != nil {
			return fmt.Errorf("Error in manga list, get %s\n", err)
		}
	default:
		return fmt.Errorf("Invalid operation")
	}
	return nil
}

func Listen(args *manga.MangaParserArgs, done chan bool) {
	for {
		select {
		case err := <-args.ErrorChan:
			args.Log.Error("Error: %s", err)
		case <-done:
			return
		}
	}
}
