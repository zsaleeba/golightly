package golightly

// ASTKind indicates which type of AST node this is
type ASTKind int

const (
	// operators
	ASTKindTopLevel ASTKind = iota
	ASTKindExpr
)

// type AST is a "sum type" implemented using an interface.
// It represents an Abstract Syntax Tree.
//
// ASTs can be created using struct initialisers.
// eg. StringToken{TokenIdentifier, "hello"}
type AST interface {
	ASTKind() ASTKind
}

type ASTTopLevel struct {
}

func (ast *ASTTopLevel) ASTKind() ASTKind {
	return ASTKindExpr
}

type ASTExpr struct {
}

func (ast *ASTExpr) ASTKind() ASTKind {
	return ASTKindExpr
}
