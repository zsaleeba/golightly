package golightly

// parseDataType parses a data type.
// if no data type is present, the first return value is false.
// Type      = TypeName | TypeLit | "(" Type ")" .
// TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
//             SliceType | MapType | ChannelType .
// TypeName  = identifier | QualifiedIdent .
func (p *Parser) parseDataType() (bool, AST, error) {
	// what token do we have?
	tok, _ := p.lexer.PeekToken(0)

	var ast AST
	var err error

	switch tok.TokenKind() {
	case TokenIdentifier:
		ast, err = p.parseOptionallyQualifiedIdentifier()

	case TokenOpenSquareBracket:
		ast, err = p.parseDataTypeArray()

	case TokenStruct:
		ast, err = p.parseDataTypeStruct()

	case TokenAsterisk:
		ast, err = p.parseDataTypePointer()

	case TokenFunc:
		ast, err = p.parseDataTypeFunction()

	case TokenInterface:
		ast, err = p.parseDataTypeInterface()

	case TokenMap:
		ast, err = p.parseDataTypeMap()

	case TokenChan:
		ast, err = p.parseDataTypeChannel()

	case TokenOpenBracket:
		ast, err = p.parseDataTypeBracketed()

	default:
		return false, nil, nil
	}

	return true, ast, err
}

// parseDataTypeArray parses an array data type or a slice data type.
// ArrayType   = "[" ArrayLength "]" ElementType .
// ArrayLength = Expression .
// ElementType = Type .
// SliceType = "[" "]" ElementType .
func (p *Parser) parseDataTypeArray() (AST, error) {
	// we already know is starts with '['
	startToken, _ := p.lexer.GetToken()

	// is the next character a ']'?
	tok, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}
	var arrayLength AST
	if tok.TokenKind() != TokenCloseSquareBracket {
		// it's an array length
		arrayLength, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	// it should be followed by a closing ']'
	endSpan, err := p.expectTokenPos(TokenCloseSquareBracket, "you need a ']' here")
	if err != nil {
		return nil, err
	}

	// now get the element type
	tok, err = p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}
	match, elementType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, NewError(p.filename, tok.Pos(), "I was looking for a data type here, but sadly I didn't get one")
	}

	// make the new data type
	var ast AST
	if arrayLength == nil {
		// it's a slice
		ast = ASTDataTypeSlice{startToken.Pos().Add(endSpan), elementType}
	} else {
		// it's an array
		ast = ASTDataTypeArray{startToken.Pos().Add(endSpan), arrayLength, elementType}
	}

	return ast, nil
}

// parseDataTypeStruct parses a struct data type.
// StructType     = "struct" "{" { FieldDecl ";" } "}" .
// FieldDecl      = (IdentifierList Type | AnonymousField) [ Tag ] .
// AnonymousField = [ "*" ] TypeName .
// Tag            = string_lit .
func (p *Parser) parseDataTypeStruct() (AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypePointer parses a pointer data type.
// PointerType = "*" BaseType .
// BaseType = Type .
func (p *Parser) parseDataTypePointer() (AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeFunction parses a function data type.
// FunctionType   = "func" Signature .
// Signature      = Parameters [ Result ] .
// Result         = Parameters | Type .
// Parameters     = "(" [ ParameterList [ "," ] ] ")" .
// ParameterList  = ParameterDecl { "," ParameterDecl } .
// ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
func (p *Parser) parseDataTypeFunction() (AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeInterface parses an interface data type.
// InterfaceType      = "interface" "{" { MethodSpec ";" } "}" .
// MethodSpec         = MethodName Signature | InterfaceTypeName .
// MethodName         = identifier .
// InterfaceTypeName  = TypeName .
func (p *Parser) parseDataTypeInterface() (AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeMap parses a map data type.
// MapType     = "map" "[" KeyType "]" ElementType .
// KeyType     = Type .
func (p *Parser) parseDataTypeMap() (AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeChannel parses a channel data type.
// ChannelType = ( "chan" [ "<-" ] | "<-" "chan" ) ElementType .
func (p *Parser) parseDataTypeChannel() (AST, error) {
	tok, _ := p.lexer.PeekToken(0)
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseDataTypeBracketed parses a data type enclosed by brackets.
func (p *Parser) parseDataTypeBracketed() (AST, error) {
	// absorb the open bracket
	p.lexer.GetToken()

	// get the data type
	tok, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	match, ast, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, NewError(p.filename, tok.Pos(), "by my reckoning this should have been a data type")
	}

	// get the close bracket
	err = p.expectToken(TokenCloseBracket, "I need a ')' here to finish the data type")
	if err != nil {
		return nil, err
	}

	return ast, err
}
