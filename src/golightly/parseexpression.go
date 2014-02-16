package golightly

// parseExpression parses an expression.
func (p *Parser) parseExpression() (AST, error) {
	tok, _ := p.lexer.GetToken()

	if tok.TokenKind() == TokenKindLiteralInt {
		return ASTValue{tok.Pos(), NewValueFromToken(tok, p.ts)}, nil
	}

	return nil, NewError(p.filename, tok.Pos(), "bad expression. bad.")
}
