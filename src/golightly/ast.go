package golightly

// type AST is a "sum type" implemented using an interface.
// It represents an Abstract Syntax Tree.
//
// ASTs can be created using struct initialisers.
// eg. StringToken{TokenIdentifier, "hello"}
type AST interface {
	IsAST()
	Pos() SrcSpan
}

type ASTTopLevel struct {
	pos           SrcSpan // where it is in the source
	packageName   string  // the name of the package everything is contained in
	imports       []AST   // import statements
	topLevelDecls []AST   // top level declarations
}

func (ast ASTTopLevel) IsAST() {
}

func (ast ASTTopLevel) Pos() SrcSpan {
	return ast.pos
}

type ASTImport struct {
	pos         SrcSpan // where the keyword is in the source
	packageName AST     // local package name to import as, or "." to import to the local scope.
	importPath  AST     // the path to the package or local package name.
}

func (ast ASTImport) IsAST() {
}

func (ast ASTImport) Pos() SrcSpan {
	return ast.pos
}

type ASTUnaryExpr struct {
	pos   SrcSpan   // where it is in the source
	op    TokenKind // what kind of operation it is
	param AST       // the parameter
}

func (ast ASTUnaryExpr) IsAST() {
}

func (ast ASTUnaryExpr) Pos() SrcSpan {
	return ast.pos
}

type ASTBinaryExpr struct {
	pos   SrcSpan   // where it is in the source
	op    TokenKind // what kind of operation it is
	left  AST       // the left parameter
	right AST       // the right parameter
}

func (ast ASTBinaryExpr) IsAST() {
}

func (ast ASTBinaryExpr) Pos() SrcSpan {
	return ast.pos
}

type ASTValue struct {
	pos SrcSpan // where it is in the source
	val Value   // the value
}

func (ast ASTValue) IsAST() {
}

func (ast ASTValue) Pos() SrcSpan {
	return ast.pos
}

func NewASTValueFromToken(v Token, ts *DataTypeStore) ASTValue {
	return ASTValue{v.Pos(), NewValueFromToken(v, ts)}
}

type ASTIdentifier struct {
	pos         SrcSpan // where it is in the source
	packageName string  // what package it's in
	name        string  // the identifier name
}

func (ast ASTIdentifier) IsAST() {
}

func (ast ASTIdentifier) Pos() SrcSpan {
	return ast.pos
}

type ASTDataType struct {
	pos SrcSpan  // where it is in the source
	typ DataType // what data type it is
}

func (ast ASTDataType) IsAST() {
}

func (ast ASTDataType) Pos() SrcSpan {
	return ast.pos
}

type ASTConstDecl struct {
	ident AST // the variable to declare
	typ   AST // the optional data type
	value AST // the value to set it to
}

func (ast ASTConstDecl) IsAST() {
}

func (ast ASTConstDecl) Pos() SrcSpan {
	return ast.ident.Pos()
}

type ASTDataTypeDecl struct {
	ident AST // the variable to declare
	typ   AST // the data type
}

func (ast ASTDataTypeDecl) IsAST() {
}

func (ast ASTDataTypeDecl) Pos() SrcSpan {
	return ast.ident.Pos()
}

type ASTDataTypeSlice struct {
	pos         SrcSpan // where the slice indicators [] are
	elementType AST     // slice of this data type
}

func (ast ASTDataTypeSlice) IsAST() {
}

func (ast ASTDataTypeSlice) Pos() SrcSpan {
	return ast.pos
}

type ASTDataTypeArray struct {
	pos         SrcSpan // where the array indicators [] are
	arraySize   AST     // how large the array is
	elementType AST     // slice of this data type
}

func (ast ASTDataTypeArray) IsAST() {
}

func (ast ASTDataTypeArray) Pos() SrcSpan {
	return ast.pos
}

type ASTDataTypePointer struct {
	pos         SrcSpan // where the pointer indicator * is
	elementType AST     // pointer to this data type
}

func (ast ASTDataTypePointer) IsAST() {
}

func (ast ASTDataTypePointer) Pos() SrcSpan {
	return ast.pos
}

type ASTDataTypeMap struct {
	pos       SrcSpan // where the map indicators map[...] are
	keyType   AST     // key is this data type
	valueType AST     // value is this data type
}

func (ast ASTDataTypeMap) IsAST() {
}

func (ast ASTDataTypeMap) Pos() SrcSpan {
	return ast.pos
}

// type ChanDirection is the directions data can travel on a channel
type ChanDirection int

const (
	ChanDirectionIn ChanDirection = iota
	ChanDirectionOut
	ChanDirectionBi
)

type ASTDataTypeChan struct {
	pos         SrcSpan       // where the chan indicators chan and <- are
	dir         ChanDirection // what directions data can flow on this channel
	elementType AST           // pointer to this data type
}

func (ast ASTDataTypeChan) IsAST() {
}

func (ast ASTDataTypeChan) Pos() SrcSpan {
	return ast.pos
}

type ASTDataTypeStruct struct {
	pos    SrcSpan // the entire struct definition
	fields []AST   // fields of this struct
}

func (ast ASTDataTypeStruct) IsAST() {
}

func (ast ASTDataTypeStruct) Pos() SrcSpan {
	return ast.pos
}

type ASTDataTypeField struct {
	identifier AST    // identifier of this field
	typ        AST    // type of this field
	tag        string // tag associated with this field
}

func (ast ASTDataTypeField) IsAST() {
}

func (ast ASTDataTypeField) Pos() SrcSpan {
	if ast.identifier != nil {
		return ast.identifier.Pos()
	} else {
		return ast.typ.Pos()
	}
}

type ASTVarDecl struct {
	ident AST // the variable to declare
	typ   AST // the optional data type
	value AST // the value to set it to
}

func (ast ASTVarDecl) IsAST() {
}

func (ast ASTVarDecl) Pos() SrcSpan {
	return ast.ident.Pos()
}
