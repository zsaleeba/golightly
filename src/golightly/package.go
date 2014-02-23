package golightly

type Package struct {
	name string  // the package name.
	syms *SymbolTable // the symbols in this package
}
