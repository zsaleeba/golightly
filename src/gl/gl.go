package main

import "fmt"

func usage() {
	fmt.Print(
`Format: goli [options] [file.go...]
	If no file arguments are provided the current directory will be searched
	for .go files.
Options:
	-s - use GoScript syntax
	-i - interactive mode
`)
}

func main() {
	fmt.Println("golightly")
	usage()
}
