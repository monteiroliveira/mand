package main

import (
	"fmt"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/pkg"
)

func main() {
	args := internal.NewArgs()
	_, err := pkg.NewParser(args)
	if err != nil {
		fmt.Printf("Error getting the parser from source, get %s\n", err)
	}
}
