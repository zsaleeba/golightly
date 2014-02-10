package golightly

import "fmt"

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

	// get a semicolon separator
	err = p.parseSemicolon("I'm gonna be needing a semicolon after this 'package' declaration")
	if err != nil {
		return err
	}

	// get a number of import declarations
	tok, err := p.lexer.PeekToken(0)
	if err != nil {
		return err
	}

	if tok.TokenKind() == TokenImport {
		for {
			// get an import
			imports, err := p.parseImport()
			if err != nil {
				return err
			}

			ast.imports = append(ast.imports, imports...)

			// get a semicolon separator
			err = p.parseSemicolon("I'm gonna be needing a semicolon after this 'import' declaration")
			if err != nil {
				return err
			}
		}
	}

	// get a number of top-level declarations
	tok, err = p.lexer.PeekToken(0)
	if err != nil {
		return err
	}

	for {
		// get a top-level declaration
		match, topLevelDecls, err := p.parseTopLevelDecl()
		if err != nil {
			return err
		}

		if !match {
			break
		}

		ast.topLevelDecls = append(ast.topLevelDecls, topLevelDecls...)

		// get a semicolon separator
		err = p.parseSemicolon("I need a semicolon here")
		if err != nil {
			return err
		}
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
func (p *Parser) parseImport() ([]AST, error) {
	// get the import declaration
	importToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}
	if importToken.TokenKind() != TokenImport {
		return nil, nil
	}

	// is it a group or a single import?
	p.lexer.GetToken()
	nextToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}
	if nextToken.TokenKind() == TokenOpenBracket {
		// get a series of import specs
		imports, err := p.parseGroup(p.parseImportSpec, "import")
		if err != nil {
			return nil, err
		}

		return imports, nil
	} else {
		// get a single import
		tree, err := p.parseImportSpec()
		if err != nil {
			return nil, err
		}

		astSlice := make([]AST, 1)
		astSlice[0] = tree
		return astSlice, nil
	}
}

// parseImportSpec parses import specifications as part of an import statement.
// ImportSpec       = [ "." | PackageName ] ImportPath .
func (p *Parser) parseImportSpec() (AST, error) {
	// what kind of thing are we looking at?
	nextToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	switch nextToken.TokenKind() {
	case TokenIdentifier:
		// it's of the form 'import fred "frod"' - get a package name first.
		strPackageName := nextToken.(StringToken)
		p.lexer.GetToken()

		// get an import path
		pathToken, err := p.lexer.GetToken()
		if err != nil {
			return nil, err
		}
		if pathToken.TokenKind() != TokenString {
			return nil, NewError(p.filename, pathToken.Pos(), "this should have been a string. eg. 'import fred \"github.com/fred/thefredpackage\"'")
		}

		return ASTImport{pathToken.Pos(), ASTIdentifier{nextToken.Pos(), strPackageName.strVal}, NewASTValueFromToken(pathToken, p.ts)}, nil

	case TokenString:
		// it's of the form 'import "frod"' - just get the import path.
		p.lexer.GetToken()
		return ASTImport{nextToken.Pos(), nil, NewASTValueFromToken(nextToken, p.ts)}, nil

	default:
		return nil, NewError(p.filename, nextToken.Pos(), "this import makes no sense. It should be like 'import [cool] \"coolpackage\"'")
	}
}

// parseTopLevelDecl parses a top-level declaration.
// TopLevelDecl  = Declaration | FunctionDecl | MethodDecl .
// Declaration   = ConstDecl | TypeDecl | VarDecl .
func (p *Parser) parseTopLevelDecl() (bool, []AST, error) {
	// what kind of thing are we looking at?
	nextToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return false, nil, err
	}

	switch nextToken.TokenKind() {
	case TokenConst:
		asts, err := p.parseConstDecl()
		return true, asts, err

	case TokenTypeKeyword:
		asts, err := p.parseTypeDecl()
		return true, asts, err

	case TokenVar:
		asts, err := p.parseVarDecl()
		return true, asts, err

	case TokenFunc:
		// is it a func decl or a method decl?
		nextToken, err = p.lexer.PeekToken(1)
		if err != nil {
			return false, nil, err
		}
		if nextToken.TokenKind() == TokenOpenBracket {
			// '(' is a total giveaway - it's a method decl
			ast, err := p.parseMethodDecl()
			return true, []AST{ast}, err
		} else {
			// it's a func decl
			ast, err := p.parseFunctionDecl()
			return true, []AST{ast}, err
		}

	default:
		return false, nil, NewError(p.filename, nextToken.Pos(), "so I wanted a top level thing like a type, a func, a const or a var, but no... you had to be different")
	}
}

// parseConstDecl parses a constant declaration.
// ConstDecl      = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
func (p *Parser) parseConstDecl() ([]AST, error) {
	// we already know it starts with "const"
	p.lexer.GetToken()

	// is it a '(' next?
	bracketToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	var decls []AST
	if bracketToken.TokenKind() == TokenOpenBracket {
		// it's a group of const specs
		decls, err = p.parseGroup(p.parseConstSpec, "const")
		if err != nil {
			return nil, err
		}
	} else {
		// it's a single const spec
		decl, err := p.parseConstSpec()
		if err != nil {
			return nil, err
		}

		decls = []AST{decl}
	}

	return decls, nil
}

// parseConstSpec parses a constant spec.
// ConstSpec      = IdentifierList [ [ Type ] "=" ExpressionList ] .
func (p *Parser) parseConstSpec() (AST, error) {
/*
	// get the identifier list
	identList, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}

	// is there a data type following?
	matchTyp, typ, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	// maybe an equals?
	var exprList []AST
	equalsToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	// handle optional part
	if matchTyp || equalsToken.TokenKind() == TokenEquals {
		// there must be an '=' and expression list after a type
		if equalsToken.TokenKind() != TokenEquals {
			return nil, NewError(p.filename, typ.Pos(), "this should really be followed by '='")
		}

		// get the expression list
		p.lexer.GetToken()
		exprList, err = p.parseExpressionList()
		if err != nil {
			return nil, err
		}
	}

	// make a const out of all this
	XXX
	 */
	// we already know it starts with "var"
	constToken, _ := p.lexer.GetToken()

	//
	return nil, NewError(p.filename, constToken.Pos(), "unimplemented")
}

// parseTypeDecl parses a type declaration.
// TypeDecl     = "type" ( TypeSpec | "(" { TypeSpec ";" } ")" ) .
func (p *Parser) parseTypeDecl() ([]AST, error) {
	// we already know it starts with "type"
	typeToken, _ := p.lexer.GetToken()

	//
	return nil, NewError(p.filename, typeToken.Pos(), "unimplemented")
}

// parseVarDecl parses a variable declaration.
// VarDecl     = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
func (p *Parser) parseVarDecl() ([]AST, error) {
	// we already know it starts with "var"
	varToken, _ := p.lexer.GetToken()

	//
	return nil, NewError(p.filename, varToken.Pos(), "unimplemented")
}

// parseFunctionDecl parses a function declaration. Note that "func" will
// already have been consumed so we're starting from the FunctionName.
// FunctionDecl = "func" FunctionName ( Function | Signature ) .
func (p *Parser) parseFunctionDecl() (AST, error) {
	// we already know it starts with "func"
	funcToken, _ := p.lexer.GetToken()

	//
	return nil, NewError(p.filename, funcToken.Pos(), "unimplemented")
}

// parseMethodDecl parses a method declaration. Note that "func" will
// already have been consumed so we're starting from the Receiver.
// MethodDecl   = "func" Receiver MethodName ( Function | Signature ) .
func (p *Parser) parseMethodDecl() (AST, error) {
	// we already know it starts with "func"
	funcToken, _ := p.lexer.GetToken()

	//
	return nil, NewError(p.filename, funcToken.Pos(), "unimplemented")
}

// parseGroup parses a group of some other clause, surrounded by brackets and
// with semicolons after each entry.
func (p *Parser) parseGroup(parseClause func()(AST, error), verbName string) ([]AST, error) {
	openBracketToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}
	if openBracketToken.TokenKind() != TokenOpenBracket {
		return nil, NewError(p.filename, openBracketToken.Pos(), "there should be a '(' here")
	}

	// get a series of sub-clauses
	p.lexer.GetToken()
	var asts []AST
	semiErrorMessage := fmt.Sprint("I really wanted a semicolon between these '", verbName, "'s")
	for {
		// is it a terminating ')'?
		closeBracketToken, err := p.lexer.PeekToken(0)
		if err != nil {
			return nil, err
		}
		if closeBracketToken.TokenKind() == TokenCloseBracket {
			break
		}

		// parse a sub-clause
		newClause, err := parseClause()
		if err != nil {
			return nil, err
		}

		// get a semicolon separator
		err = p.parseSemicolon(semiErrorMessage)
		if err != nil {
			return nil, err
		}

		asts = append(asts, newClause)
	}

	return asts, nil
}

// parseSemicolon parses a required semicolon
func (p *Parser) parseSemicolon(message string) error {
	// get a semicolon separator
	semicolonToken, err := p.lexer.GetToken()
	if err != nil {
		return err
	}
	if semicolonToken.TokenKind() != TokenSemicolon {
		return NewError(p.filename, semicolonToken.Pos(), message)
	}

	return nil
}
