package golightly

// parseExpressionList parses a comma-separated list of expressions.
// ExpressionList = Expression { "," Expression } .
func (p *Parser) parseExpressionList() ([]AST, error) {
	// get an expression
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	asts := []AST{expr}

	// get more commas then expressions
	for {
		// look for a comma
		comma, err := p.lexer.PeekToken(0)
		if err != nil {
			return nil, err
		}

		if comma.TokenKind() != TokenKindComma {
			break
		}

		p.lexer.GetToken()

		// get an expression
		expr, err = p.parseExpression()
		if err != nil {
			return nil, err
		}

		// add the identifier to our list of identifiers
		asts = append(asts, expr)

	}

	return asts, nil
}

// parseExpression parses an expression.
func (p *Parser) parseExpression() (AST, error) {
	tok, _ := p.lexer.GetToken()

	if tok.TokenKind() == TokenKindLiteralInt {
		return ASTValue{tok.Pos(), NewValueFromToken(tok, p.ts)}, nil
	}

	return nil, NewError(p.filename, tok.Pos(), "bad expression. bad.")
}
