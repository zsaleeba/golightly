package golightly

// parseStatement parses a statement.
func (p *Parser) parseStatement() (AST, error) {
	tok, _ := p.lexer.GetToken()
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}

// parseBlock parses a statement block
// Block = "{" StatementList "}" .
// StatementList = { Statement ";" } .
func (p *Parser) parseBlock() (AST, error) {
	tok, _ := p.lexer.GetToken()
	return nil, NewError(p.filename, tok.Pos(), "unimplemented")
}
