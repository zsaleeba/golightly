package main

import (
	"fmt"
	"golightly"
	"os"
	"runtime"
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

	// allow it to use all the CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// create the compiler
	c := golightly.NewCompiler()

	// compile the program
	err := c.Compile(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
