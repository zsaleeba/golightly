package golightly

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const (
	addImportChannelDepth  = 4
	compileSrcChannelDepth = 4
	completionChannelDepth = 4
)

// type compileStatus
type compileStatus int
const (
	compileStatusParsing = iota
	compileStatusSymbolsAvailable
	compileStatusComplete
)

// type Compiler is a container for an entire compiler, including
// all passes of all files.
//
// The compiler class has two main tasks:
//   * compiling files
//   * importing packages
//
// It tries to do this while allowing for a reasonable amount of concurrent
// compilation. All Go files are scheduled for lexing and parsing
// concurrently. On a multi-core machine this may result in parallel
// compilation.
//
// LEXING AND PARSING
//
// The lexer is called directly from the parser so lexing and parsing
// occur as a single pass. This pass takes source code as input and
// outputs an abstract syntax tree (AST).
//
// When imports are parsed the packages are scheduled for concurrent
// importation. The symbols from the imports aren't needed until after
// parsing is complete so imports can occur concurrently with parsing.
//
// If a package is already compiled imports may only need to read the
// pre-compiled package. Alternatively it may require that the package
// be compiled first. In that case more Go files from the package will
// be scheduled for compilation.
//
// After each Go file is parsed it's necessary to ensure that all the
// import dependencies have finished importing before continuing to
// the next phase. At this point the file's goroutine waits on all
// the imports completing before it continues.
//
// Once a file has completed parsing it signals this to whatever
// code requested its compilation. If it's a file from an imported
// package this lets the package know that it's ok to use the symbols.
// When all files in the package have been parsed the package can
// signal that its symbols are ready to any source files which have
// imported it.
//
// When all the files which were requested from the command line have
// finished parsing the compiler can proceed to the next pass -
// semantic analysis.
//
// SEMANTIC ANALYSIS
//
// Semantic analysis is a multi-pass process which checks a lot of
// semantics of the program and transforms the AST into a somewhat
// simpler form which is easier for subsequent passes to work with.
// Symbols are resolved in this pass.
//
// Semantic analysis starts with package main, function main. All
// symbols referenced in function main are resolved and then resolution
// is performed on those functions, and so on until all the symbols
// which are used have been processed - and none of the symbols
// which are unused are processed.
//
// Each time a new symbol is analysed it's checked for two things:
//    - whether the file it's in has changed since the previous
//      compilation, and
//    - whether the AST for this symbol is identical to the previous
//      compilation.
//
// The AST checksum from the previous compilation is stored in a
// database for comparison purposes. Unless the symbol is changed due
// to either of the above circumstances all of the following passes
// will be omitted and the symbol will retrieve its target executable
// code from the database and go straight to linking.
//
// AST OPTIMISATION
//
// AST optimisation performs a series of optimisations on the AST,
// looking for patterns which can be transformed into more efficient
// forms.
//
// IR PROCESSING
//
// IR processing transforms the AST into a DAG intermediate
// representation and then performs a series of optimisations on
// the IR.
//
// CODE GENERATION
//
// Code generation transforms the IR into target executable code
// plus debug and link information.
//
// LINK
//
//
//
type Compiler struct {
	srcFiles map[string]*sourceFile    // the files we're compiling.
	packages map[string]*compilePackage // the packages we're importing or defining.

	shutdown chan bool // closed when the compiler is shutting down.

	dataTypeStore *DataTypeStore // keeps a global set of data types known to the compiler.

	addImport  chan importMessage     // new packages are queued for import using this stream.
	compileSrc chan compileSrcMessage // new files are queued for compilation using this stream.
}

// type importMessage is sent to Compiler.addImport to request that a package be imported.
type importMessage struct {
	packageName     string                 // the requested package name to import.
	fromFileName    string                 // what source file it was requested from.
	pos             SrcSpan                // where in the source file it was requested from.
	completeChannel chan completionMessage // how to notify when it's done.
}

// type compileSrcMessage is sent to Compiler.compileSrc to request that a file be compiled.
type compileSrcMessage struct {
	fileName        string
	completeChannel chan completionMessage
}

// type completionMessage is sent to notify a caller of completion of
// compilation, with a possible error.
type completionMessage struct {
	packageName string // what package we were working on.
	fileName    string // what file we were working on.
	err         error  // error from compilation or nil on success.
}

// NewCompiler creates a new compiler object.
func NewCompiler() *Compiler {
	c := new(Compiler)

	c.srcFiles = make(map[string]*sourceFile)
	c.packages = make(map[string]*compilePackage)

	c.shutdown = make(chan bool)

	c.dataTypeStore = NewDataTypeStore()
	c.addImport = make(chan importMessage, addImportChannelDepth)
	c.compileSrc = make(chan compileSrcMessage, compileSrcChannelDepth)

	// accept source files for compilation
	go c.parseSrcs()

	// accept packages to import
	go c.importPackages()

	return c
}

func (c *Compiler) Close() {
}

// Compile is the central point to compile a program from. It takes
// all the files as arguments and produces a runnable program as
// output. All passes of the compiler are run.
func (c *Compiler) Compile(srcFiles []string) error {
	// create a channel for source files to notify us when their symbols are ready.
	completeChannel := make(chan completionMessage, completionChannelDepth)

	// queue the source files for compilation.
	// these are picked up by parseSrcs() and compiled.
	waitingOn := make(map[string]bool)
	for _, fileName := range srcFiles {
		// are we already compiling it?
		_, found := waitingOn[fileName]
		if !found {
			// need to compile it.
			waitingOn[fileName] = true
			c.compileSrc <- compileSrcMessage{fileName, completeChannel}
		}
	}

	// wait for symbols ready or error.
	var err error
	for {
		// get a message from a compilation.
		msg := <-completeChannel

		// either got "symbols ready" from a file or an error.
		if msg.err != nil {
			err = msg.err
			close(c.shutdown) // tell it to shutdown.
		}

		delete(waitingOn, msg.fileName)
		if len(waitingOn) == 0 {
			// we've finished all of them.
			break
		}
	}

	return err
}

// parseFileAndComplete parses a single file, called from schedulePass. To compile a file
// you should send it to the Compiler.compileSrc channel for parseSrcs() to
// compile. After the file is parsed a completion message is sent to the client.
func (c *Compiler) parseFileAndComplete(sf *sourceFile) {
	err := c.parseFile(sf)
	if err != nil {
		sf.completeChannel <- completionMessage{sf.packageName, sf.fileName, )}
		return
	}
}


// compileFile parses a single file, called from schedulePass. To compile a file
// you should send it to the Compiler.compileSrc channel for parseSrcs() to
// compile.
func (c *Compiler) compileFile(sf *sourceFile) error {
	// open the source file
	srcFile, err := os.Open(sf.fileName)
	if err != nil {
		return errors.New(fmt.Sprintf("I can't find ", sf.fileName, ": ", err)
	}

	defer srcFile.Close()
	srcReader := bufio.NewReader(srcFile)

	// lex and parse it.
	lex := NewLexer()
	lex.LexReader(srcReader, sf.fileName)
	parser := NewParser(lex, c.dataTypeStore, sf)
	err = parser.Parse()
	if err != nil {
		return err
	}

	// create symbols.
	err = c.createSymbols(sf)
	if err != nil {
		return err
	}

	// wait for imports to complete.
	err = c.waitImports(sf)
	if err != nil {
		return err
	}

	// say we're done.
	return nil
}

// createSymbols creates a set of symbols from an already parsed source file.
// when we're finished we tell our parent package that we're done.
func (c *Compiler) createSymbols(sf *sourceFile) error {
	return nil
}

// compileSrcs runs as a goroutine, accepting files to parse and
// parsing them.
func (c *Compiler) compileSrcs() {
	for {
		// wait for something to happen.
		var running bool

		select {
		case csm := <-c.compileSrc:
			// add to srcFiles.
			sf := NewSourceFile(csm.fileName, c.compileSrc, c.addImport, csm.completeChannel, c.shutdown)
			c.srcFiles[csm.fileName] = sf

			// start parsing the file
			go c.compileFileAndComplete(sf)

		case _, running = <-c.shutdown:
			// running is false if we're shutting down.
		}

		// are we shutting down?
		if !running {
			break
		}
	}
}

// importPackages runs as a goroutine, accepting packages to import and
// importing them.
func (c *Compiler) importPackages() {
	importComplete := make(chan completionMessage, completionChannelDepth)

	for {
		// wait for something to happen.
		var running bool

		select {
		case im := <-c.addImport:
			// a new package to import. do we already know about it?
			cp, ok := c.packages[im.packageName]
			if ok {
				// we're already importing this package.
				if cp.status == compileStatusParsing {
					// add to the list of clients to be informed when it's done.
					cp.clientCompleteChannels = append(cp.clientCompleteChannels, im.completeChannel)
				} else {
					// let the client know immediate that we're done.
					im.completeChannel <- cp.completeMessage
				}
			} else {
				// add to packages.
				cp = NewCompilePackage(im.packageName, c.compileSrc, c.addImport, importComplete, c.shutdown)
				c.packages[im.packageName] = cp
			}

		case cm := <-importComplete:
			// we got a completion message from a package.
			cp, ok := c.packages[cm.packageName]
			if ok {
				// keep the completion message in case we need it for a later import.
				cp.completeMessage = cm

				// tell everyone who wants to know.
				for _, client := range cp.clientCompleteChannels {
					client <- cm
				}
				cp.clientCompleteChannels = nil
				cp.status = compileStatusSymbolsAvailable
			}

		case _, running = <-c.shutdown:
			// running is false if we're shutting down.
		}

		// are we shutting down?
		if !running {
			break
		}
	}
}
