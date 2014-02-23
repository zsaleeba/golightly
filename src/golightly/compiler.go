package golightly

import (
	"errors"
	"os"
	"fmt"
	"bufio"
	"runtime"
)

// type compilerPass indicates a stage in compilation which a source file
// can be at.
type compilerPass int
const (
	compilerPassStart = iota
	compilerPassParse
	compilerPassResolve
	compilerPassSemantic
	compilerPassCodeGen
	compilerPassLink
)

const (
	addImportChannelDepth = 4
	compileSrcChannelDepth = 4
	progressChannelDepth = 4
)

// type Compiler is a container for an entire compiler, including
// all passes of all files.
type Compiler struct {
	readyToCompile      []*SourceFile // each of the source files to be compiled.
	compiling map[string]*SourceFile // each of the source files currently being compiled.

	maxThreads int  // the maximum number of threads to use in compilation.
	currentThreads int // the current number of compilation threads running.
	shutdown bool      // if true the compiler is shutting down. we won't schedule
					   // any more work.

	dataTypeStore *DataTypeStore // keeps a global set of data types known to the compiler.

	errorStream chan error  // errors are sent to this channel for display.
	addImport   chan importMessage  // new packages are queued for import using this stream.
	compileSrc   chan string  // new files are queued for compilation using this stream.
	progress chan string // progress in compilation is sent here.
	finished chan bool // used to inform Compile() that all the threads are done.
}

// type importMessage is sent to Compiler.addImport to request that a package be imported.
type importMessage struct {
	packageName string  // the requested package name to import
	fromFileName string // what source file it was requested from
	pos SrcSpan // where in the source file it was requested from
}

// NewCompiler creates a new compiler object.
func NewCompiler() *Compiler {
	c := new(Compiler)

	c.readyToCompile = nil
	c.compiling = make(map[string]*SourceFile)

	c.maxThreads = runtime.NumCPU()
	c.currentThreads = 0
	c.shutdown = false

	c.dataTypeStore = NewDataTypeStore()
	c.addImport = make(chan importMessage, addImportChannelDepth)
	c.compileSrc = make(chan string, compileSrcChannelDepth)
	c.progress = make(chan string, progressChannelDepth)
	c.errorStream = make(chan error)
	c.finished = make(chan bool)

	// accept source files for compilation
	go c.compileSrcs()

	// accept packages to import
	go c.importPackages()

	return c
}

func (c *Compiler) Close() {
}

// SetMaxThreads is used before compilation to set the maximum number of
// threads which will be used in compilation.
func (c *Compiler) SetMaxThreads(maxThreads int) {
	c.maxThreads = maxThreads
}

// Compile is the central point to compile a program from. It takes
// all the files as arguments and produces a runnable program as
// output. All passes of the compiler are run.
func (c *Compiler) Compile(srcFiles []string) error {
	// queue the source files for compilation.
	// these are picked up by compileSrcs() and compiled.
	for _, fileName := range srcFiles {
		c.compileSrc <- fileName
	}

	// wait for an error message or completion.
	var err error
	select {
	case err = <- c.errorStream:
		// got an error - shutdown first.
		c.shutdown = true  // tell it to shutdown.
		<- c.finished      // wait until it's done.

	case <- c.finished:
		// got successful completion.
	}

	return err
}

// compileSrcs runs as a goroutine, accepting files to compile and compiling
// them. It runs each compilation pass on each of the source files.
func (c *Compiler) compileSrcs() {
	for {
		// wait for something to happen.
		select {
		case newFileName := <- c.compileSrc:
			// a new source file. queue it to compile.
			newSrcFile := new(SourceFile)
			newSrcFile.fileName = newFileName
			newSrcFile.pass = compilerPassStart
			c.readyToCompile = append(c.readyToCompile, newSrcFile)

		case progFileName := <- c.progress:
			// we've finished a pass while compiling something.
			srcFile := c.getFromCompiling(progFileName)
			if srcFile.pass != compilerPassLink {
				// we've finished a pass - queue it for the next pass.
				c.readyToCompile = append(c.readyToCompile, srcFile)
			}
		}

		// have we aborted due to an error?
		if c.shutdown {
			break
		}

		// schedule something new to compile.
		done := c.schedulePass()
		if done {
			break
		}
	}

	// shutting down, just wait until each of our threads is done.
	for {
		// are we done?
		if len(c.compiling) == 0 {
			break
		}

		// wait for something to finish
		select {
		case <- c.compileSrc:
			// a new file? don't care.

		case progFileName := <- c.progress:
			// a pass finished. maybe we're done.
			c.getFromCompiling(progFileName)
		}
	}

	// tell Compile() we're done.
	c.finished <- true
}

// getFromCompiling gets a named file from the set of "currently
// being compiled" files.
func (c *Compiler) getFromCompiling(fileName string) *SourceFile {
	sf, ok := c.compiling[fileName]
	if !ok {
		return nil
	}

	// remove it from the map
	delete(c.compiling, fileName)

	return sf
}

// schedulePass takes a source file from "readyToCompile", puts it in the
// "compiling" map and starts it compiling.
func (c *Compiler) schedulePass() bool {
	// do this as many times as we can
	for ; c.currentThreads < c.maxThreads && len(c.readyToCompile) > 0; {
		// get the first thing off readyToCompile.
		sf := c.readyToCompile[0]
		c.readyToCompile = c.readyToCompile[1:]

		// put it in "compiling".
		c.compiling[sf.fileName] = sf

		// start the next pass running.
		switch sf.pass {
		case compilerPassStart:
			sf.pass = compilerPassParse
			go c.parseFile(sf)

		case compilerPassParse:
			sf.pass = compilerPassResolve
			go c.resolveSymbols(sf)

		case compilerPassResolve:
			sf.pass = compilerPassSemantic
			go c.semanticAnalysis(sf)

		case compilerPassSemantic:
			sf.pass = compilerPassCodeGen
			go c.codeGeneration(sf)

		case compilerPassLink:
			sf.pass = compilerPassLink
			go c.link(sf)
		}
	}

	// are we done?
	done := len(c.readyToCompile) == 0 && len(c.compiling) == 0
	return done
}

// parseFile parses a single file, called from schedulePass. To compile a file
// you should send it to the Compiler.compileSrc channel for compileSrcs() to
// compile.
func (c *Compiler) parseFile(sf *SourceFile) {
	// make sure we tell compileSrcs() we're done before we return
	fileName := sf.fileName
	defer func() {
		c.progress <- fileName
	}()

	// open the source file
	srcFile, err := os.Open(fileName)
	if err != nil {
		c.errorStream <- errors.New(fmt.Sprintf("I can't find ", fileName, ": ", err))
		return
	}

	defer srcFile.Close()
	srcReader := bufio.NewReader(srcFile)

	// lex and parse it
	lex := NewLexer()
	lex.LexReader(srcReader, fileName)
	parser := NewParser(lex, c.dataTypeStore, c.compileSrc)
	err = parser.Parse()

	if err != nil {
		c.errorStream <- err
	}
}

// resolveSymbols
func (c *Compiler) resolveSymbols(sf *SourceFile) {
	// make sure we tell compileSrcs() we're done before we return
	fileName := sf.fileName
	defer func() {
		c.progress <- fileName
	}()
}

// semanticAnalysis
func (c *Compiler) semanticAnalysis(sf *SourceFile) {
	// make sure we tell compileSrcs() we're done before we return
	fileName := sf.fileName
	defer func() {
		c.progress <- fileName
	}()
}

// codeGeneration
func (c *Compiler) codeGeneration(sf *SourceFile) {
	// make sure we tell compileSrcs() we're done before we return
	fileName := sf.fileName
	defer func() {
		c.progress <- fileName
	}()
}

// link
func (c *Compiler) link(sf *SourceFile) {
	// make sure we tell compileSrcs() we're done before we return
	fileName := sf.fileName
	defer func() {
		c.progress <- fileName
	}()
}



// importPackages runs as a goroutine, accepting packages to import and
// importing them.
func (c *Compiler) importPackages() {
	for {
		// wait for something to happen.
		select {
		case <- c.addImport:
			// a new package to import.

		}

		// have we aborted due to an error?
		if c.shutdown {
			break
		}
	}
}
