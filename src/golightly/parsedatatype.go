package golightly

// parseDataType parses a data type.
// if no data type is present, the first return value is false.
// Type      = TypeName | TypeLit | "(" Type ")" .
// TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
//             SliceType | MapType | ChannelType .
func (p *Parser) parseDataType() (bool, AST, error) {
	// what token do we have?
	tok, _ := p.lexer.PeekToken(0)

	switch tok.TokenKind() {
	case TokenIdentifier:
		return p.parseDataTypeName()

	case TokenOpenSquareBracket:
		return p.parseDataTypeArray()

	case TokenStruct:
		return p.parseDataTypeStruct()

	case TokenAsterisk:
		return p.parseDataTypePointer()

	case TokenFunc:
		return p.parseDataTypeFunction()

	case TokenInterface:
		return p.parseDataTypeInterface()

	case TokenMap:
		return p.parseDataTypeMap()

	case TokenChan:
		return p.parseDataTypeChannel()

	case TokenOpenBracket:
		p.lexer.GetToken()
		match, ast, err := p.parseDataType()
		if err != nil {
			return match, nil, err
		}

		err = p.expectToken(TokenCloseBracket, "I need a ')' here to finish the data type")
		if err != nil {
			return match, nil, err
		}

		return match, ast, nil

	default:
		return false, nil, nil
	}
}

// parseDataType parses a data type name.
// TypeName  = identifier | QualifiedIdent .
func (p *Parser) parseDataTypeName() (bool, AST, error) {
	ast, err := p.parseOptionallyQualifiedIdentifier()
	if err != nil {
		return false, nil, err
	}

	return true, ast, nil
}

// parseDataTypeArray parses an array data type.
// ArrayType   = "[" ArrayLength "]" ElementType .
// ArrayLength = Expression .
// ElementType = Type .
// SliceType = "[" "]" ElementType .
func (p *Parser) parseDataTypeArray() (bool, AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return false, nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeStruct parses a struct data type.
// StructType     = "struct" "{" { FieldDecl ";" } "}" .
// FieldDecl      = (IdentifierList Type | AnonymousField) [ Tag ] .
// AnonymousField = [ "*" ] TypeName .
// Tag            = string_lit .
func (p *Parser) parseDataTypeStruct() (bool, AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return false, nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypePointer parses a pointer data type.
// PointerType = "*" BaseType .
// BaseType = Type .
func (p *Parser) parseDataTypePointer() (bool, AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return false, nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeFunction parses a function data type.
// FunctionType   = "func" Signature .
// Signature      = Parameters [ Result ] .
// Result         = Parameters | Type .
// Parameters     = "(" [ ParameterList [ "," ] ] ")" .
// ParameterList  = ParameterDecl { "," ParameterDecl } .
// ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
func (p *Parser) parseDataTypeFunction() (bool, AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return false, nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeInterface parses an interface data type.
// InterfaceType      = "interface" "{" { MethodSpec ";" } "}" .
// MethodSpec         = MethodName Signature | InterfaceTypeName .
// MethodName         = identifier .
// InterfaceTypeName  = TypeName .
func (p *Parser) parseDataTypeInterface() (bool, AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return false, nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeMap parses a map data type.
// MapType     = "map" "[" KeyType "]" ElementType .
// KeyType     = Type .
func (p *Parser) parseDataTypeMap() (bool, AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return false, nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeChannel parses a channel data type.
// ChannelType = ( "chan" [ "<-" ] | "<-" "chan" ) ElementType .
func (p *Parser) parseDataTypeChannel() (bool, AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return false, nil, NewError(p.filename, tok.Pos(), "unimplemented")
}
