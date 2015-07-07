package golightly

// type compilePackage is a package which is imported or defined by the source code.
type compilePackage struct {
	packageName         string                   // the name of this package.
	symbols             SymbolTable              // the symbols in this package - only valid once symbol creation is complete for all package files.
	waitingFileComplete map[string]bool          // the files from this package we're still waiting on.
	fileComplete        chan completionMessage   // files tell us they're complete with a message on this channel.
	compileSrc          chan compileSrcMessage   // we can request files to be compiled here.
	addImport           chan importMessage       // we can request imports here.
	completeChannel     chan completionMessage   // channel to importPackages() to notify when our symbols are complete.
	shutdown            chan bool                // closed when the compiler is shutting down.

	// the following are used by Compiler.importPackages().
	status				compileStatus            // where we are in the compilation process.
	clientCompleteChannels    []chan completionMessage // channels back to clients for importPackages() to notify when our symbols are complete.
	completeMessage     completionMessage        // importPackages() uses this internally.
}

// NewCompilePackage creates a new compilePackage.
func NewCompilePackage(packageName string, compileSrc chan compileSrcMessage, addImport chan importMessage, completeChannel chan completionMessage, shutdown chan bool) *compilePackage {
	sp := new(compilePackage)
	sp.packageName = packageName
	sp.waitingFileComplete = make(map[string]bool)
	sp.fileComplete = make(chan completionMessage)
	sp.compileSrc = compileSrc
	sp.addImport = addImport
	sp.completeChannel = make(chan completionMessage, completionChannelDepth)
	sp.shutdown = shutdown

	sp.clientCompleteChannels = make([]chan completionMessage, 1)
	sp.clientCompleteChannels[0] = completeChannel

	return sp
}
