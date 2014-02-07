package golightly

// type Parser controls parsing of a token stream into an AST
type Parser struct {
	filename    string // the name of the file being parsed
	packageName string // the name of the package this file is a part of

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

// parseSourceFile parses the contents of an entire source file
func (p *Parser) parseSourceFile() error {
	// get the package declaration
	err := p.parsePackage()
	if err != nil {
		return err
	}

	// get a number of import declarations
	for {
		match, err := p.parseImport()
		if err != nil {
			return err
		}

		if !match {
			break
		}
	}

	// get a number of top-level declarations
	for {
		match, err := p.parseTopLevelDecl()
		if err != nil {
			return err
		}

		if !match {
			break
		}
	}

	// make sure we're at the end of the file
	endToken := p.GetToken(0, 0)
	if endToken.TokenKind() != TokenEndOfSource {
		return NewError(p.filename, endToken.Pos(), "I don't really know what this is or why it's here")
	}

	return nil
}

// parseSourceFile parses a package declaration
func (p *Parser) parsePackage() error {
	// get the package declaration
	packageToken := p.GetToken(0, 1)
	if packageToken.TokenKind() != TokenPackage {
		return NewError(p.filename, packageToken.Pos(), "the file should start with 'package <package name>'")
	}

	packageNameToken := p.GetToken(0, 1)
	if packageNameToken.TokenKind() != TokenIdentifier {
		return NewError(p.filename, packageToken.Pos(), "the package name should be a plain word. eg. 'package horatio'")
	}

	strPackageName := packageNameToken.(StringToken)
	p.packageName = strPackageName.strVal

	return nil
}

// parseImport parses an import declaration
func (p *Parser) parseImport() (bool, error) {
	// get the import declaration
	importToken := p.GetToken(0, 0)
	if importToken.TokenKind() != TokenImport {
		return false, nil
	}

	// is it a group or a single import?
	nextToken := p.GetToken(1, 1)
	if nextToken.TokenKind() == TokenOpenGroup {
		// get a series of quoted strings with semicolons
	} else {
		// get a single quoted string
		if nextToken.TokenKind() != TokenString {
			return true, NewError(p.filename, nextToken.Pos(), "imports should be a quoted string. eg. 'import \"cubancigars\"'")
		}

		//strPackageName := nextToken.(StringToken)
		//XXX = strPackageName.strVal
		//work out how to import stuff!
	}

	return false, nil
}

// parseTopLevelDecl parses a top-level declaration
func (p *Parser) parseTopLevelDecl() (bool, error) {
	return false, nil
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
