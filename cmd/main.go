package main

import (
	"fmt"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg"
)

func main() {
	args := internal.NewArgs()
	parser, err := pkg.NewParser(args)
	if err != nil {
		fmt.Printf("Error getting the parser from source, get %s\n", err)
	}
	pages, err := parser.ExtractSingleChapter()
	if err != nil {
		fmt.Printf("Error in source content extraction, get %s\n", err)
	}
	err = parser.DownloadPages(pages)
	if err != nil {
		fmt.Printf("Error in pages download, get %s\n", err)
	}
}
