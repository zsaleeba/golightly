package golightly

type SourceFile struct {
	pkg *Package  // which package this file is in
	fileName string // the name of the file
	pass compilerPass // what compiler pass we've most recently completed
	ast AST  // the parse tree from this file
}
