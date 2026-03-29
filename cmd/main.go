package main

import (
	"errors"
	"fmt"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg"
	"github.com/monteiroliveira/mand/pkg/parsers/manga"
)

func main() {
	args := internal.NewArgs()

	// TODO: create a mand module and initialize cli in internal (cli) package
	switch {
	case args.Manga != nil:
		mangaArgs := &manga.MangaParserArgs{
			Verbose: args.Verbose,
		}
		if args.Manga.Download != nil {
			mangaArgs.Operation = manga.DownloadOperation
			mangaArgs.Source = args.Manga.Download.Source.URL
		}
		if args.Manga.DownloadList != nil {
			mangaArgs.Operation = manga.DownloadListOperation
			mangaArgs.Source = args.Manga.DownloadList.Source.URL
		}

		parser, err := pkg.NewMangaParser(mangaArgs)
		if err != nil {
			fmt.Printf("Error getting the parser from source, get %s\n", err)
		}
		if err = pkg.Execute(parser, mangaArgs.Operation); err != nil {
			err = errors.Join(fmt.Errorf("Error executing command"), err)
			fmt.Println(err)
		}
	}
}
