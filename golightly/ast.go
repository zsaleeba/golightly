package golightly

// type AST is a "sum type" implemented using an interface.
// It represents an Abstract Syntax Tree.
//
// ASTs can be created using struct initialisers.
// eg. StringToken{TokenIdentifier, "hello"}
type AST interface {
	IsAST()
	Pos() SrcSpan
	Equals(to AST) bool
}

// type ASTTopLevel describes the top level of a source file.
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

func (ast ASTTopLevel) Equals(to AST) bool {
	too := to.(ASTTopLevel)
	if !ast.pos.Equals(too.pos) ||
		ast.packageName != too.packageName ||
		len(ast.imports) != len(too.imports) ||
		len(ast.topLevelDecls) != len(too.topLevelDecls) {
		return false
	}

	for i, a := range ast.imports {
		if !a.Equals(too.imports[i]) {
			return false
		}
	}

	for i, a := range ast.topLevelDecls {
		if !a.Equals(too.topLevelDecls[i]) {
			return false
		}
	}

	return true
}

// type ASTImport describes an import statement.
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

func (ast ASTImport) Equals(to AST) bool {
	too := to.(ASTImport)
	return ast.pos.Equals(too.pos) && ast.packageName.Equals(too.packageName) && ast.importPath.Equals(too.importPath)
}

// type ASTUnaryExpr describes an expression operation with a single operand.
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

func (ast ASTUnaryExpr) Equals(to AST) bool {
	too := to.(ASTUnaryExpr)
	return ast.pos.Equals(too.pos) && ast.op == too.op && ast.param.Equals(too.param)
}

// type ASTBinaryExpr describes an expression operation with two operands.
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

func (ast ASTBinaryExpr) Equals(to AST) bool {
	too := to.(ASTBinaryExpr)
	return ast.pos.Equals(too.pos) && ast.op == too.op && ast.left.Equals(too.left) && ast.right.Equals(too.right)
}

// type ASTValue describes a literal value.
type ASTValue struct {
	pos SrcSpan // where it is in the source
	val Value   // the value
}

func (ast ASTValue) IsAST() {
}

func (ast ASTValue) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTValue) Equals(to AST) bool {
	too := to.(ASTValue)
	return ast.pos.Equals(too.pos) && ast.val.Equals(too.val)
}

func NewASTValueFromToken(v Token, ts *DataTypeStore) ASTValue {
	return ASTValue{v.Pos(), NewValueFromToken(v, ts)}
}

// type ASTIdentifier describes an identifier reference.
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

func (ast ASTIdentifier) Equals(to AST) bool {
	too := to.(ASTIdentifier)
	return ast.pos.Equals(too.pos) && ast.packageName == too.packageName && ast.name == too.name
}

// type ASTConstDecl describes a constant declaration.
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

func (ast ASTConstDecl) Equals(to AST) bool {
	too := to.(ASTConstDecl)
	return ast.ident.Equals(too.ident) && ast.typ.Equals(too.typ) && ast.value.Equals(too.value)
}

// type ASTVarDecl describes a variable declaration.
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

func (ast ASTVarDecl) Equals(to AST) bool {
	too := to.(ASTVarDecl)
	return ast.ident.Equals(too.ident) && ast.typ.Equals(too.typ) && ast.value.Equals(too.value)
}

// type ASTFunctionDecl describes a function or method declaration.
type ASTFunctionDecl struct {
	pos      SrcSpan // the 'func <name>' part of the declaration
	name     string  // the function name
	receiver AST     // the optional receiver
	params   []AST   // the parameters
	returns  []AST   // the return values
	body     AST     // the body of the function
}

func (ast ASTFunctionDecl) IsAST() {
}

func (ast ASTFunctionDecl) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTFunctionDecl) Equals(to AST) bool {
	too := to.(ASTFunctionDecl)
	if !(ast.pos.Equals(too.pos) && ast.name == too.name && ast.receiver.Equals(too.receiver) && ast.body.Equals(too.body)) {
		return false
	}

	if len(ast.params) != len(too.params) || len(ast.returns) != len(too.returns) {
		return false
	}

	for i, a := range ast.params {
		if !a.Equals(too.params[i]) {
			return false
		}
	}

	for i, a := range ast.returns {
		if !a.Equals(too.returns[i]) {
			return false
		}
	}

	return true
}

// type ASTReceiver describes a receiver in a method declaration.
type ASTReceiver struct {
	pos      SrcSpan // the whole receiver
	name     string  // the receiving variable name
	pointer  bool    // true if it's of the form *Type
	typeName string  // the name of the receiver's type
}

func (ast ASTReceiver) IsAST() {
}

func (ast ASTReceiver) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTReceiver) Equals(to AST) bool {
	too := to.(ASTReceiver)
	return ast.pos.Equals(too.pos) && ast.name == too.name && ast.pointer == too.pointer && ast.typeName == too.typeName
}

// type ASTDataTypeDecl describes a type declaration using the 'type' keyword.
type ASTDataTypeDecl struct {
	ident AST // the variable to declare
	typ   AST // the data type
}

func (ast ASTDataTypeDecl) IsAST() {
}

func (ast ASTDataTypeDecl) Pos() SrcSpan {
	return ast.ident.Pos()
}

func (ast ASTDataTypeDecl) Equals(to AST) bool {
	too := to.(ASTDataTypeDecl)
	return ast.ident.Equals(too.ident) && ast.typ.Equals(too.typ)
}

// type ASTDataTypeSlice describes a slice declaration.
type ASTDataTypeSlice struct {
	pos         SrcSpan // where the slice indicators [] are
	elementType AST     // slice of this data type
}

func (ast ASTDataTypeSlice) IsAST() {
}

func (ast ASTDataTypeSlice) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTDataTypeSlice) Equals(to AST) bool {
	too := to.(ASTDataTypeSlice)
	return ast.pos.Equals(too.pos) && ast.elementType.Equals(too.elementType)
}

// type ASTDataTypeArray describes an array declaration.
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

func (ast ASTDataTypeArray) Equals(to AST) bool {
	too := to.(ASTDataTypeArray)
	return ast.pos.Equals(too.pos) && ast.arraySize.Equals(too.arraySize) && ast.elementType.Equals(too.elementType)
}

// type ASTDataTypePointer describes a pointer declaration.
type ASTDataTypePointer struct {
	pos         SrcSpan // where the pointer indicator * is
	elementType AST     // pointer to this data type
}

func (ast ASTDataTypePointer) IsAST() {
}

func (ast ASTDataTypePointer) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTDataTypePointer) Equals(to AST) bool {
	too := to.(ASTDataTypePointer)
	return ast.pos.Equals(too.pos) && ast.elementType.Equals(too.elementType)
}

// type ASTDataTypeMap describes a map declaration.
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

func (ast ASTDataTypeMap) Equals(to AST) bool {
	too := to.(ASTDataTypeMap)
	return ast.pos.Equals(too.pos) && ast.keyType.Equals(too.keyType) && ast.valueType.Equals(too.valueType)
}

// type ChanDirection is the directions data can travel on a channel.
type ChanDirection int

const (
	ChanDirectionIn ChanDirection = iota
	ChanDirectionOut
	ChanDirectionBi
)

// type ASTDataTypeChan describes a channel declaration.
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

func (ast ASTDataTypeChan) Equals(to AST) bool {
	too := to.(ASTDataTypeChan)
	return ast.pos.Equals(too.pos) && ast.dir == too.dir && ast.elementType.Equals(too.elementType)
}

// type ASTDataTypeStruct describes a structure declaration.
type ASTDataTypeStruct struct {
	pos    SrcSpan // the entire struct definition
	fields []AST   // fields of this struct
}

func (ast ASTDataTypeStruct) IsAST() {
}

func (ast ASTDataTypeStruct) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTDataTypeStruct) Equals(to AST) bool {
	too := to.(ASTDataTypeStruct)
	if !ast.pos.Equals(too.pos) && len(ast.fields) == len(too.fields) {
		return false
	}

	for i, a := range ast.fields {
		if !a.Equals(too.fields[i]) {
			return false
		}
	}

	return true
}

// type ASTDataTypeField describes a field of a struct.
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

func (ast ASTDataTypeField) Equals(to AST) bool {
	too := to.(ASTDataTypeField)
	return ast.identifier.Equals(too.identifier) && ast.typ.Equals(too.typ) && ast.tag == too.tag
}

// type ASTDataTypeFunc describes a function/method declaration.
type ASTDataTypeFunc struct {
	pos     SrcSpan // the entire func signature
	params  []AST   // parameters
	returns []AST   // return values of this function
}

func (ast ASTDataTypeFunc) IsAST() {
}

func (ast ASTDataTypeFunc) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTDataTypeFunc) Equals(to AST) bool {
	too := to.(ASTDataTypeFunc)
	if !(ast.pos.Equals(too.pos) && len(ast.params) == len(too.params) || len(ast.returns) == len(too.returns)) {
		return false
	}

	for i, a := range ast.params {
		if !a.Equals(too.params[i]) {
			return false
		}
	}

	for i, a := range ast.returns {
		if !a.Equals(too.returns[i]) {
			return false
		}
	}

	return true
}

// type ASTParamDecl describes a function/method parameter or return value.
type ASTParameterDecl struct {
	identifier AST // the name of the parameter
	typ        AST // the type of the parameter
}

func (ast ASTParameterDecl) IsAST() {
}

func (ast ASTParameterDecl) Pos() SrcSpan {
	if ast.identifier != nil {
		return ast.identifier.Pos().Add(ast.typ.Pos())
	} else {
		return ast.typ.Pos()
	}
}

func (ast ASTParameterDecl) Equals(to AST) bool {
	too := to.(ASTParameterDecl)
	return ast.identifier.Equals(too.identifier) && ast.typ.Equals(too.typ)
}

// type ASTEllipsis describes an ellipsis as part of a parameter list.
type ASTEllipsis struct {
	pos SrcSpan // where the ellipsis is
}

func (ast ASTEllipsis) IsAST() {
}

func (ast ASTEllipsis) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTEllipsis) Equals(to AST) bool {
	too := to.(ASTEllipsis)
	return ast.pos.Equals(too.pos)
}

// type ASTDataTypeInterface describes an interface declaration.
type ASTDataTypeInterface struct {
	pos     SrcSpan // the start of the interface definition
	methods []AST   // methods of this interface
}

func (ast ASTDataTypeInterface) IsAST() {
}

func (ast ASTDataTypeInterface) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTDataTypeInterface) Equals(to AST) bool {
	too := to.(ASTDataTypeInterface)
	if !(ast.pos.Equals(too.pos) && len(ast.methods) == len(too.methods)) {
		return false
	}

	for i, a := range ast.methods {
		if !a.Equals(too.methods[i]) {
			return false
		}
	}

	return true
}

// type ASTDataTypeMethodSpec describes a method within an interface declaration.
type ASTDataTypeMethodSpec struct {
	pos     SrcSpan // where the name is in the source
	name    string  // the identifier name
	params  []AST   // the method parameters
	returns []AST   // the method return values
}

func (ast ASTDataTypeMethodSpec) IsAST() {
}

func (ast ASTDataTypeMethodSpec) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTDataTypeMethodSpec) Equals(to AST) bool {
	too := to.(ASTDataTypeMethodSpec)
	if !(ast.pos.Equals(too.pos) && ast.name == too.name && len(ast.params) == len(too.params) || len(ast.returns) == len(too.returns)) {
		return false
	}

	for i, a := range ast.params {
		if !a.Equals(too.params[i]) {
			return false
		}
	}

	for i, a := range ast.returns {
		if !a.Equals(too.returns[i]) {
			return false
		}
	}

	return true
}

// type ASTBlock describes a block and the statements in it.
type ASTBlock struct {
	pos        SrcSpan // the entire span of the block
	statements []AST   // the statements in the block
}

func (ast ASTBlock) IsAST() {
}

func (ast ASTBlock) Pos() SrcSpan {
	return ast.pos
}

func (ast ASTBlock) Equals(to AST) bool {
	too := to.(ASTBlock)
	if !(ast.pos.Equals(too.pos) && len(ast.statements) == len(too.statements)) {
		return false
	}

	for i, a := range ast.statements {
		if !a.Equals(too.statements[i]) {
			return false
		}
	}

	return true
}
