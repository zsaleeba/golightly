package golightly

import (
	"errors"
	"sync"
)

// type Compiler is a container for an entire compiler, including
// all passes of all files.
type Compiler struct {
	files      map[string]SourceFile
	filesMutex sync.Mutex
}

// NewCompiler creates a new compiler object.
func NewCompiler() *Compiler {
	c := new(Compiler)
	return c
}

// Compile is the central point to compile a program from. It takes
// all the files as arguments and produces a runnable program as
// output. All passes of the compiler are run.
func (c *Compiler) Compile(srcFiles []string) error {
	return errors.New("unimplemented")
}
