package golightly

// parseStatement parses a statement.
// Statement =
// Declaration | LabeledStmt | SimpleStmt |
// GoStmt | ReturnStmt | BreakStmt | ContinueStmt | GotoStmt |
// FallthroughStmt | Block | IfStmt | SwitchStmt | SelectStmt | ForStmt |
// DeferStmt .
// SimpleStmt = EmptyStmt | ExpressionStmt | SendStmt | IncDecStmt | Assignment | ShortVarDecl .
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
