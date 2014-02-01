package golightly

import "errors"

// type Parser controls parsing of a token stream into an AST
type Parser struct {
	filename    string // the name of the file being parsed

	in chan Token // the stream of tokens is received through this channel
	tokenQueue  []Token // a buffer of the incoming tokens

	out chan *AST // the stream of ASTs is sent out through this channel
}

func NewParser() *Parser {
	p := new(Parser)
	p.out = make(chan *AST)
	p.tokenQueue = make([]Token, 0, 3)
	return p
}

// Parse runs the parser and breaks the program down into an Abstract Syntax Tree.
func (p *Parser) Parse() error {
	return nil
}

func (p *Parser) parseSourceFile() error {
	packageToken := p.GetToken(0, 1)
	if packageToken.TokenKind() != TokenPackage {
		return NewError(p.filename, packageToken.Pos(), "the file should start with 'package <package name>'")
	}

	return errors.New("unimplemented")
}

// GetToken gets a token from the input token channel, with look-ahead available.
func (p *Parser) GetToken(ahead int, discard int) Token {
	// do we need to get more tokens?
	if ahead >= len(p.tokenQueue) {
		// do we need to make the token queue larger?
		if ahead >= cap(p.tokenQueue) {
			// make more space
			newQueue := make([]Token, len(p.tokenQueue), ahead+1)
			copy(newQueue, p.tokenQueue)
			p.tokenQueue = newQueue
		}

		// get some more tokens
		for len(p.tokenQueue) <= ahead {
			newToken := <- p.in
			p.tokenQueue = append(p.tokenQueue, newToken)
		}
	}

	// discard some if we have to
	result := p.tokenQueue[ahead]
	if discard > 0 {
		copy(p.tokenQueue, p.tokenQueue[discard:])
	}

	return result
}
