package main

import (
	"fmt"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg"
)

func main() {
	args := internal.NewArgs()

	// TODO: create a mand module and initialize cli in internal (cli) package
	switch {
	case args.Manga != nil:
		if args.Manga.Download != nil {
			parser, err := pkg.NewMangaParser(args.Manga.Download.Source.URL)
			if err != nil {
				fmt.Printf("Error getting the parser from source, get %s\n", err)
			}
			pages, err := parser.ExtractSingleChapter()
			if err != nil {
				fmt.Printf("Error in source content extraction, get %s\n", err)
			}
			chn, err := parser.ExtractChapterName()
			if err != nil {
				fmt.Printf("Error in get chapter name, get %s\n", err)
			}
			err = parser.DownloadPages(pages, chn)
			if err != nil {
				fmt.Printf("Error in pages download, get %s\n", err)
			}
		}
		if args.Manga.DownloadList != nil {
			parser, err := pkg.NewMangaParser(args.Manga.DownloadList.Source.URL)
			if err != nil {
				fmt.Printf("Error getting the parser from source, get %s\n", err)
			}
			err = parser.ExtractChapterList()
			if err != nil {
				fmt.Printf("Error in manga list, get %s\n", err)
			}
		}
	}
}
