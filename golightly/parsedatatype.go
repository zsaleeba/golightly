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
	case TokenKindIdentifier:
		ast, err = p.parseOptionallyQualifiedIdentifier()

	case TokenKindOpenSquareBracket:
		ast, err = p.parseDataTypeArray()

	case TokenKindStruct:
		ast, err = p.parseDataTypeStruct()

	case TokenKindAsterisk:
		ast, err = p.parseDataTypePointer()

	case TokenKindFunc:
		ast, err = p.parseDataTypeFunction()

	case TokenKindInterface:
		ast, err = p.parseDataTypeInterface()

	case TokenKindMap:
		ast, err = p.parseDataTypeMap()

	case TokenKindChan, TokenKindChannelArrow:
		ast, err = p.parseDataTypeChannel()

	case TokenKindOpenBracket:
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
	if tok.TokenKind() != TokenKindCloseSquareBracket {
		// it's an array length
		arrayLength, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	// it should be followed by a closing ']'
	endSpan, err := p.expectTokenPos(TokenKindCloseSquareBracket, "you need a ']' here")
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
		if arrayLength == nil {
			return nil, NewError(p.filename, tok.Pos(), "I was looking for a data type in this slice definition - it should look like '[]element_type'")
		} else {
			return nil, NewError(p.filename, tok.Pos(), "I was looking for a data type in this array definition - it should look like '[size]element_type'")
		}
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
func (p *Parser) parseDataTypeStruct() (AST, error) {
	// get the 'struct' token
	structTok, _ := p.lexer.GetToken()

	// get a '{' as well
	err := p.expectToken(TokenKindOpenBrace, "struct definitions need a '{' here")
	if err != nil {
		return nil, err
	}

	// get the struct fields
	var fields []AST
	for {
		// are we at the end?
		tok, err := p.lexer.PeekToken(0)
		if err != nil {
			return nil, err
		}

		if tok.TokenKind() == TokenKindCloseBrace {
			break
		}

		// get a field
		newFields, err := p.parseDataTypeField()
		if err != nil {
			return nil, err
		}

		fields = append(fields, newFields...)

		// get a semicolon
		err = p.expectToken(TokenKindSemicolon, "semicolon expected between struct fields")
		if err != nil {
			return nil, err
		}
	}

	// get the trailing '}'
	endPos, err := p.expectTokenPos(TokenKindCloseBrace, "struct definitions need a '}' here")
	if err != nil {
		return nil, err
	}

	return ASTDataTypeStruct{structTok.Pos().Add(endPos), fields}, nil
}

// parseDataTypeField parses a struct field declaration.
// FieldDecl      = (IdentifierList Type | AnonymousField) [ Tag ] .
// AnonymousField = [ "*" ] TypeName .
// Tag            = string_lit .
func (p *Parser) parseDataTypeField() ([]AST, error) {
	// what do we have here?
	tok, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	var idents []AST
	if tok.TokenKind() == TokenKindIdentifier {
		// try parsing it as an identifier list
		idents, err = p.parseIdentifierList("struct field")
		if err != nil {
			return nil, err
		}
	}

	// what type were these identifiers?
	typeTok, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	match, typ, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, NewError(p.filename, typeTok.Pos(), "I needed a data type here in this struct field declaration")
	}

	// get a trailing tag if one exists
	var tag string
	tagTok, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}
	if tagTok.TokenKind() == TokenKindLiteralString {
		tag = tagTok.(StringToken).strVal
		p.lexer.GetToken()
	}

	// make the result
	if idents == nil {
		// just return a single anonymous field
		return []AST{ASTDataTypeField{nil, typ, tag}}, nil
	} else {
		// return a set of struct fields
		fields := make([]AST, len(idents))
		for i, ident := range idents {
			fields[i] = ASTDataTypeField{ident, typ, tag}
		}

		return fields, nil
	}
}

// parseDataTypePointer parses a pointer data type.
// PointerType = "*" BaseType .
// BaseType = Type .
func (p *Parser) parseDataTypePointer() (AST, error) {
	// get the '*' token
	tok, _ := p.lexer.GetToken()

	// get the element type
	tok2, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	match, elementType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, NewError(p.filename, tok2.Pos(), "by my reckoning this part of a pointer definition should have been a data type")
	}

	return ASTDataTypePointer{tok.Pos(), elementType}, nil
}

// parseDataTypeFunction parses a function data type.
// FunctionType   = "func" Signature .
// Signature      = Parameters [ Result ] .
// Result         = Parameters | Type .
func (p *Parser) parseDataTypeFunction() (AST, error) {
	// get the "func" token
	funcTok, _ := p.lexer.GetToken()

	// get a function signature
	params, returns, err := p.parseSignature()
	if err != nil {
		return nil, err
	}

	return ASTDataTypeFunc{funcTok.Pos(), params, returns}, nil
}

// parseDataTypeInterface parses an interface data type.
// InterfaceType      = "interface" "{" { MethodSpec ";" } "}" .
// MethodSpec         = MethodName Signature | InterfaceTypeName .
// MethodName         = identifier .
// InterfaceTypeName  = TypeName .
func (p *Parser) parseDataTypeInterface() (AST, error) {
	// get the 'interface' token
	interfaceToken, _ := p.lexer.GetToken()

	// get a '{' as well
	err := p.expectToken(TokenKindOpenBrace, "interface definitions need a '{' here")
	if err != nil {
		return nil, err
	}

	// get the interface methods
	var methods []AST
	for {
		// are we at the end?
		tok, err := p.lexer.PeekToken(0)
		if err != nil {
			return nil, err
		}

		if tok.TokenKind() == TokenKindCloseBrace {
			break
		}

		// get a field
		method, err := p.parseDataTypeMethodSpec()
		if err != nil {
			return nil, err
		}

		methods = append(methods, method)

		// get a semicolon
		err = p.expectToken(TokenKindSemicolon, "semicolon expected between interface methods")
		if err != nil {
			return nil, err
		}
	}

	// get the trailing '}'
	err = p.expectToken(TokenKindCloseBrace, "interface definitions need a '}' here")
	if err != nil {
		return nil, err
	}

	return ASTDataTypeInterface{interfaceToken.Pos(), methods}, nil
}

// parseDataTypeMethodSpec parses an interface data type.
// MethodSpec         = MethodName Signature | InterfaceTypeName .
// MethodName         = identifier .
// InterfaceTypeName  = TypeName .
func (p *Parser) parseDataTypeMethodSpec() (AST, error) {
	// if it's a method name the second token will be '(' to start the signature.
	tok2, err := p.lexer.PeekToken(1)
	if err != nil {
		return nil, err
	}

	if tok2.TokenKind() == TokenKindOpenBracket {
		// it's a method name
		methodName, err := p.lexer.GetToken()
		if err != nil {
			return nil, err
		}

		if methodName.TokenKind() != TokenKindIdentifier {
			return nil, NewError(p.filename, methodName.Pos(), "this should be a method name, but I'm not really seeing it")
		}

		// get the signature
		params, returns, err := p.parseSignature()
		if err != nil {
			return nil, err
		}

		return ASTDataTypeMethodSpec{methodName.Pos(), methodName.(StringToken).strVal, params, returns}, nil
	} else {
		// it must be an interface type name
		return p.parseOptionallyQualifiedIdentifier()
	}
}

// parseDataTypeMap parses a map data type.
// MapType     = "map" "[" KeyType "]" ElementType .
// KeyType     = Type .
func (p *Parser) parseDataTypeMap() (AST, error) {
	// get the 'map' token
	mapToken, _ := p.lexer.GetToken()

	// get the opening '['
	openSquareBracketToken, err := p.lexer.GetToken()
	if err != nil {
		return nil, err
	}
	if openSquareBracketToken.TokenKind() == TokenKindOpenSquareBracket {
		return nil, NewError(p.filename, mapToken.Pos().Add(openSquareBracketToken.Pos()), "map types should look like 'map[key_type]element_type'")
	}

	// get the key type
	tok, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	match, keyType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, NewError(p.filename, tok.Pos(), "by my reckoning this part of a map definition should have been a data type. map types should look like 'map[key_type]element_type'")
	}

	// get the closing ']'
	closeSquareBracketToken, err := p.lexer.GetToken()
	if err != nil {
		return nil, err
	}
	if closeSquareBracketToken.TokenKind() == TokenKindCloseSquareBracket {
		return nil, NewError(p.filename, closeSquareBracketToken.Pos(), "map types should look like 'map[key_type]element_type'")
	}

	// get the element type
	match, elementType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, NewError(p.filename, closeSquareBracketToken.Pos(), "by my reckoning this should have been followed by a data type. map types should look like 'map[key_type]element_type'")
	}

	return ASTDataTypeMap{mapToken.Pos().Add(closeSquareBracketToken.Pos()), keyType, elementType}, nil
}

// parseDataTypeChannel parses a channel data type.
// ChannelType = ( "chan" [ "<-" ] | "<-" "chan" ) ElementType .
func (p *Parser) parseDataTypeChannel() (AST, error) {
	var dir ChanDirection
	tok, _ := p.lexer.GetToken()
	chanSpan := tok.Pos()
	if tok.TokenKind() == TokenKindChan {
		// starts with "chan", what's next?
		tok2, err := p.lexer.PeekToken(0)
		if err != nil {
			return nil, err
		}

		if tok2.TokenKind() == TokenKindChannelArrow {
			// it's 'chan <-'
			dir = ChanDirectionIn
			chanSpan.end = tok2.Pos().end
			p.lexer.GetToken()
		}
	} else {
		// starts with '<-', we need a 'chan' now
		p.lexer.GetToken()
		tok2pos, err := p.expectTokenPos(TokenKindChan, "channels should look like 'chan', '<- chan' or 'chan <-'")
		if err != nil {
			return nil, err
		}

		chanSpan.end = tok2pos.end
	}

	// get the element type
	tok, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	match, elementType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, NewError(p.filename, tok.Pos(), "by my reckoning this part of a chan definition should have been a data type")
	}

	return ASTDataTypeChan{chanSpan, dir, elementType}, nil
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
	err = p.expectToken(TokenKindCloseBracket, "I need a ')' here to finish the data type")
	if err != nil {
		return nil, err
	}

	return ast, err
}
