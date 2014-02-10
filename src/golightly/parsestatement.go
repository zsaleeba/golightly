package golightly

// parseStatement parses a statement.
func (p *Parser) parseStatement() (AST, error) {
	tok, _ := p.lexer.GetToken()
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}
