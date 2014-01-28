package main

import (
	"fmt"
	"golightly"
	"os"
)

func usage() {
	fmt.Print(
`Format: gl [options] [<file.go>|<directory>]...
	If no file arguments are provided the current directory will be
	searched for .go files.

Options:
	-s - use GoScript syntax
	-i - interactive mode
`)
}

func main() {
	fmt.Println("golightly")

	c := golightly.NewCompiler()
	err := c.Compile(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
