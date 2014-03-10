package golightly

// type sourceFile is a single file which has to be compiled.
type sourceFile struct {
	packageName            string                 // the package name of this file.
	fileName               string                 // the name of this file. unique system-wide.
	ast                    AST                    // the AST result of parsing.
	symbols                SymbolTable            // the symbols in this file.
	waitingPackageComplete map[string]bool        // the import packages we're waiting on before we can do symbol resolution.
	packageComplete        chan completionMessage // packages tell us they're complete with a message on this channel.
	compileSrc             chan compileSrcMessage // we can request files to be compiled here.
	addImport              chan importMessage     // we can request imports here.
	completeChannel        chan completionMessage // a channel to notify when our symbols are complete.
	shutdown               chan bool              // closed when the compiler is shutting down.

	// the following are used by Compiler.parseSrcs().
	status				compileStatus            // where we are in the compilation process.
}

// NewSourceFile creates a new sourceFile.
func NewSourceFile(fileName string, compileSrc chan compileSrcMessage, addImport chan importMessage, completeChannel chan completionMessage, shutdown chan bool) *sourceFile {
	sf := new(sourceFile)
	sf.fileName = fileName
	sf.waitingPackageComplete = make(map[string]bool)
	sf.packageComplete = make(chan completionMessage)
	sf.addImport = addImport
	sf.completeChannel = completeChannel
	sf.shutdown = shutdown

	return sf
}
