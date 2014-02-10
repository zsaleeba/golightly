package golightly

// parseExpression parses an expression.
func (p *Parser) parseExpression() (AST, error) {
	tok, _ := p.lexer.GetToken()
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}
