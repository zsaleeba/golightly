package golightly

// parseDataType parses a data type.
// if no data type is present, the first return value is false.
func (p *Parser) parseDataType() (bool, AST, error) {
	tok, _ := p.lexer.PeekToken(0)

	if tok.TokenKind() == TokenInt {
		p.lexer.GetToken()
		return true, ASTDataType{tok.Pos(), p.ts.IntType()}, nil
	}

	return false, nil, nil
}
