package golightly

// type Parser controls parsing of a token stream into an AST
type Parser struct {
	lexer *Lexer // the lexical analyser
	ts *DataTypeStore // the data type store

	filename    string // the name of the file being parsed
	packageName string // the name of the package this file is a part of
}

// NewParser
func NewParser(lexer *Lexer, ts *DataTypeStore) *Parser {
	p := new(Parser)
	p.lexer = lexer
	p.ts = ts

	return p
}

// Parse runs the parser and breaks the program down into an Abstract Syntax Tree.
func (p *Parser) Parse() error {
	return nil
}

// parseSourceFile parses the contents of an entire source file.
// SourceFile       = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .
func (p *Parser) parseSourceFile() error {
	// get the package declaration
	ast := new(ASTTopLevel)
	packageName, err := p.parsePackage()
	if err != nil {
		return err
	}
	ast.packageName = packageName

	// get a number of import declarations
	for {
		match, imports, err := p.parseImport()
		if err != nil {
			return err
		}

		if !match {
			break
		}

		ast.imports = append(ast.imports, imports...)
	}

	// get a number of top-level declarations
	for {
		match, topLevelDecls, err := p.parseTopLevelDecl()
		if err != nil {
			return err
		}

		if !match {
			break
		}

		ast.topLevelDecls = append(ast.topLevelDecls, topLevelDecls...)
	}

	// make sure we're at the end of the file
	endToken, err := p.lexer.GetToken()
	if err != nil {
		return err
	}
	if endToken.TokenKind() != TokenEndOfSource {
		return NewError(p.filename, endToken.Pos(), "I don't really know what this is or why it's here")
	}

	return nil
}

// parsePackage parses a package declaration.
// PackageClause  = "package" PackageName .
func (p *Parser) parsePackage() (string, error) {
	// get the package declaration
	packageToken, err := p.lexer.GetToken()
	if err != nil {
		return "", err
	}
	if packageToken.TokenKind() != TokenPackage {
		return "", NewError(p.filename, packageToken.Pos(), "the file should start with 'package <package name>'")
	}

	packageNameToken, err := p.lexer.GetToken()
	if err != nil {
		return "", err
	}
	if packageNameToken.TokenKind() != TokenIdentifier {
		return "", NewError(p.filename, packageToken.Pos(), "the package name should be a plain word. eg. 'package horatio'")
	}

	strPackageName := packageNameToken.(StringToken)

	return strPackageName.strVal, nil
}

// parseImport parses an import declaration.
// ImportDecl       = "import" ( ImportSpec | "(" { ImportSpec ";" } ")" ) .
func (p *Parser) parseImport() (bool, []AST, error) {
	// get the import declaration
	importToken, err := p.lexer.PeekToken()
	if err != nil {
		return false, nil, err
	}
	if importToken.TokenKind() != TokenImport {
		return false, nil, nil
	}

	// is it a group or a single import?
	p.lexer.GetToken()
	nextToken, err := p.lexer.PeekToken()
	if err != nil {
		return false, nil, err
	}
	if nextToken.TokenKind() == TokenOpenGroup {
		// get a series of import specs
		p.lexer.GetToken()
		imports := make([]AST, 0, 20)
		for {
			// get an import spec
			match, newImport, err := p.parseImportSpec(importToken)
			if err != nil {
				return false, nil, err
			}
			if !match {
				break
			}

			// get a semicolon separate
			semicolonToken, err := p.lexer.GetToken()
			if err != nil {
				return false, nil, err
			}
			if semicolonToken.TokenKind() != TokenSemicolon {
				return false, nil, NewError(p.filename, semicolonToken.Pos(), "I really wanted a semicolon between these imports")
			}

			imports = append(imports, newImport)
		}

		// check for a trailing ')'
		closeBracketToken, err := p.lexer.GetToken()
		if err != nil {
			return false, nil, err
		}
		if closeBracketToken.TokenKind() != TokenCloseGroup {
			return false, nil, NewError(p.filename, closeBracketToken.Pos(), "I was looking for a ')' to finish off this import")
		}

		return true, imports, nil
	} else {
		// get a single import
		match, tree, err := p.parseImportSpec(importToken)
		if err != nil {
			return false, nil, err
		}
		if !match {
			return false, nil, NewError(p.filename, nextToken.Pos(), "this import makes no sense. It should be like 'import [cool] \"coolpackage\"'")
		}

		astSlice := make([]AST, 1)
		astSlice[0] = tree
		return true, astSlice, nil
	}

	return false, nil, nil
}

// parseImportSpec parses import specifications as part of an import statement.
// ImportSpec       = [ "." | PackageName ] ImportPath .
func (p *Parser) parseImportSpec(importToken Token) (bool, AST, error) {
	nextToken, err := p.lexer.PeekToken()
	if err != nil {
		return false, nil, err
	}

	switch nextToken.TokenKind() {
	case TokenIdentifier:
		// it's of the form 'import fred "frod"' - get a package name first.
		strPackageName := nextToken.(StringToken)
		p.lexer.GetToken()

		// get an import path
		pathToken, err := p.lexer.GetToken()
		if err != nil {
			return false, nil, err
		}
		if pathToken.TokenKind() != TokenString {
			return false, nil, NewError(p.filename, pathToken.Pos(), "this should have been a string. eg. 'import fred \"github.com/fred/thefredpackage\"'")
		}

		return true, ASTImport{importToken.Pos(), ASTIdentifier{nextToken.Pos(), strPackageName.strVal}, NewASTValueFromToken(pathToken, p.ts)}, nil

	case TokenString:
		// it's of the form 'import "frod"' - just get the import path.
		p.lexer.GetToken()
		return true, ASTImport{importToken.Pos(), nil, NewASTValueFromToken(nextToken, p.ts)}, nil

	default:
		return false, nil, NewError(p.filename, nextToken.Pos(), "this import makes no sense. It should be like 'import [cool] \"coolpackage\"'")
	}
}

// parseTopLevelDecl parses a top-level declaration
func (p *Parser) parseTopLevelDecl() (bool, []AST, error) {
	return false, nil, nil
}

/*
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
*/
