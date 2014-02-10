package golightly

// parseDataType parses a data type.
// if no data type is present, the first return value is false.
func (p *Parser) parseDataType() (bool, AST, error) {
	tok, _ := p.lexer.GetToken()
	return false, nil, NewError(p.filename, tok.Pos(), "unimplemented")
}
