package main

import (
	"errors"
	"fmt"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg"
	"github.com/monteiroliveira/mand/pkg/parsers/manga"
)

func main() {
	done := make(chan bool)
	args := internal.NewArgs()

	log := internal.NewLogger(args.Verbose)

	// TODO: create a mand module and initialize cli in internal (cli) package
	switch {
	case args.Manga != nil:
		mangaArgs := &manga.MangaParserArgs{
			Log:       log,
			ErrorChan: make(chan error),
		}
		if args.Manga.Download != nil {
			mangaArgs.Operation = manga.DownloadOperation
			mangaArgs.Source = args.Manga.Download.Source.URL
		}
		if args.Manga.DownloadList != nil {
			mangaArgs.Operation = manga.DownloadListOperation
			mangaArgs.Source = args.Manga.DownloadList.Source.URL
			mangaArgs.ListBatchSize = args.Manga.DownloadList.BatchListSize // Default value for now
		}

		parser, err := pkg.NewMangaParser(mangaArgs)
		if err != nil {
			log.Error("Error getting the parser from source, get %s", err)
			return
		}
		go pkg.Listen(mangaArgs, done)
		if err = pkg.Execute(parser, mangaArgs); err != nil {
			err = errors.Join(fmt.Errorf("Error executing command"), err)
			log.Error("%s", err)
		}
		done <- true
	}
}
