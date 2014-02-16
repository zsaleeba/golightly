package golightly

import (
	"fmt"
)

// type Parser controls parsing of a token stream into an AST
type Parser struct {
	lexer *Lexer         // the lexical analyser
	ts    *DataTypeStore // the data type store

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
	err = p.expectToken(TokenSemicolon, "I'm gonna be needing a semicolon after this 'package' declaration")
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
			err = p.expectToken(TokenSemicolon, "I'm gonna be needing a semicolon after this 'import' declaration")
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
		err = p.expectToken(TokenSemicolon, "I need a semicolon here")
		if err != nil {
			return err
		}
	}

	// make sure we're at the end of the file
	err = p.expectToken(TokenEndOfSource, "I don't really know what this is or why it's here")
	if err != nil {
		return err
	}

	return nil
}

// parsePackage parses a package declaration.
// PackageClause  = "package" PackageName .
func (p *Parser) parsePackage() (string, error) {
	// get the package declaration
	err := p.expectToken(TokenPackage, "the file should start with 'package <package name>'")
	if err != nil {
		return "", err
	}

	packageNameToken, err := p.lexer.GetToken()
	if err != nil {
		return "", err
	}
	if packageNameToken.TokenKind() != TokenIdentifier {
		return "", NewError(p.filename, packageNameToken.Pos(), "the package name should be a plain word. eg. 'package horatio'")
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
		imports, err := p.parseGroupSingle(p.parseImportSpec, "import")
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

		return ASTImport{pathToken.Pos(), ASTIdentifier{nextToken.Pos(), "", strPackageName.strVal}, NewASTValueFromToken(pathToken, p.ts)}, nil

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
		asts, err := p.parseDecl(p.parseConstSpec, "const")
		return true, asts, err

	case TokenTypeKeyword:
		asts, err := p.parseDecl(p.parseTypeSpec, "type")
		return true, asts, err

	case TokenVar:
		asts, err := p.parseDecl(p.parseVarSpec, "var")
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

// parseDecl parses a declaration. It's used for const, type and var
// declarations since they're all fairly similar.
// ConstDecl      = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
// TypeDecl       = "type"  ( TypeSpec  | "(" { TypeSpec  ";" } ")" ) .
// VarDecl        = "var"   ( VarSpec   | "(" { VarSpec   ";" } ")" ) .
func (p *Parser) parseDecl(parseSpec func() ([]AST, error), verbName string) ([]AST, error) {
	// we already know it starts with the verb, so skip that
	p.lexer.GetToken()

	// is it a '(' next?
	bracketToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	var decls []AST
	if bracketToken.TokenKind() == TokenOpenBracket {
		// it's a group of specs
		decls, err = p.parseGroupMulti(parseSpec, verbName)
		if err != nil {
			return nil, err
		}
	} else {
		// it's a single spec
		decls, err = parseSpec()
		if err != nil {
			return nil, err
		}
	}

	return decls, nil
}

// parseConstSpec parses a constant spec.
// ConstSpec      = IdentifierList [ [ Type ] "=" ExpressionList ] .
func (p *Parser) parseConstSpec() ([]AST, error) {
	// get the identifier list
	identList, err := p.parseIdentifierList("constant")
	if err != nil {
		return nil, err
	}

	// is there a data type following?
	matchTyp, typeAST, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	// maybe an equals?
	equalsToken, err := p.lexer.PeekToken(0)
	if err != nil {
		return nil, err
	}

	// handle optional part
	var exprList []AST
	if matchTyp || equalsToken.TokenKind() == TokenEquals {
		// there must be an '=' and expression list after a type
		if equalsToken.TokenKind() != TokenEquals {
			return nil, NewError(p.filename, equalsToken.Pos(), "after a data type I expected to see '=' here")
		}

		// get the expression list
		p.lexer.GetToken()
		exprList, err = p.parseExpressionList()
		if err != nil {
			return nil, err
		}
	}

	// are the two lists the same length?
	identSpan := identList[0].Pos().Add(identList[len(identList)-1].Pos())
	if len(identList) > len(exprList) {
		return nil, NewError(p.filename, identSpan, "there are more names here than there are values")
	} else if len(identList) < len(exprList) {
		return nil, NewError(p.filename, identSpan, "there are less names here than there are values")
	}

	// make a set of consts out of all this
	asts := make([]AST, len(identList))
	for i := 0; i < len(identList); i++ {
		asts[i] = ASTConstDecl{identList[i], typeAST, exprList[i]}
	}

	return asts, nil
}

// parseTypeSpec parses a type declaration specification.
// TypeSpec     = identifier Type .
func (p *Parser) parseTypeSpec() ([]AST, error) {
	// get an identifier
	ident, err := p.lexer.GetToken()
	if err != nil {
		return nil, err
	}

	if ident.TokenKind() != TokenIdentifier {
		return nil, NewError(p.filename, ident.Pos(), fmt.Sprint("this should have been a name for a type, but it's not"))
	}

	identAST := ASTIdentifier{ident.Pos(), "", ident.(StringToken).strVal}

	// get the data type
	matchTyp, typeAST, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	// the type is mandatory here
	if !matchTyp {
		fail, err := p.lexer.PeekToken(0)
		if err != nil {
			return nil, err
		}

		return nil, NewError(p.filename, fail.Pos(), fmt.Sprint("this should have been a name for a type, but it's not"))
	}

	return []AST{ASTDataTypeDecl{identAST, typeAST}}, nil
}

// parseVarSpec parses a variable declaration specification.
// VarSpec     = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
func (p *Parser) parseVarSpec() ([]AST, error) {
	// get the identifier list
	identList, err := p.parseIdentifierList("variable")
	if err != nil {
		return nil, err
	}

	// is there a data type following?
	matchTyp, typeAST, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	var exprList []AST
	if matchTyp {
		// optional equals
		equalsToken, err := p.lexer.PeekToken(0)
		if err != nil {
			return nil, err
		}

		if equalsToken.TokenKind() == TokenEquals {
			// get the expression list
			p.lexer.GetToken()
			exprList, err = p.parseExpressionList()
			if err != nil {
				return nil, err
			}
		}
	} else {
		// required equals
		err := p.expectToken(TokenEquals, "I was expecting to see an '=' here")
		if err != nil {
			return nil, err
		}

		// get the expression list
		p.lexer.GetToken()
		exprList, err = p.parseExpressionList()
		if err != nil {
			return nil, err
		}
	}

	// are the two lists the same length?
	if exprList != nil {
		identSpan := identList[0].Pos().Add(identList[len(identList)-1].Pos())

		if len(identList) > len(exprList) {
			return nil, NewError(p.filename, identSpan, "there are more names here than there are values")
		} else if len(identList) < len(exprList) {
			return nil, NewError(p.filename, identSpan, "there are less names here than there are values")
		}
	}

	// make a set of variable declarations out of all this
	asts := make([]AST, len(identList))
	for i := 0; i < len(identList); i++ {
		asts[i] = ASTVarDecl{identList[i], typeAST, exprList[i]}
	}

	return asts, nil
}

// parseIdentifierList parses a comma-separated list of identifiers.
// IdentifierList = identifier { "," identifier } .
func (p *Parser) parseIdentifierList(identDesc string) ([]AST, error) {
	var asts []AST

	for {
		// get an identifier
		ident, err := p.lexer.GetToken()
		if err != nil {
			return nil, err
		}

		if ident.TokenKind() != TokenIdentifier {
			return nil, NewError(p.filename, ident.Pos(), fmt.Sprint("this should have been a name for a ", identDesc, ", but it's not"))
		}

		// add the identifier to our list of identifiers
		asts = append(asts, ASTIdentifier{ident.Pos(), "", ident.(StringToken).strVal})

		// look for a comma after it
		comma, err := p.lexer.PeekToken(0)
		if err != nil {
			return nil, err
		}

		if comma.TokenKind() != TokenComma {
			break
		}

		p.lexer.GetToken()
	}

	return asts, nil
}

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

		if comma.TokenKind() != TokenComma {
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

// parseFunctionDecl parses a function declaration. Note that "func" will
// already have been consumed so we're starting from the FunctionName.
// FunctionDecl = "func" FunctionName ( Function | Signature ) .
func (p *Parser) parseFunctionDecl() (AST, error) {
	// we already know it starts with "func"
	funcToken, _ := p.lexer.GetToken()

	// get an identifier for the function name
	ident, err := p.lexer.GetToken()
	if err != nil {
		return nil, err
	}

	if ident.TokenKind() != TokenIdentifier {
		return nil, NewError(p.filename, ident.Pos(), fmt.Sprint("this should have been a function name, but it's not"))
	}

	// get a function or a signature

	return nil, NewError(p.filename, funcToken.Pos(), "unimplemented")
}

// parseMethodDecl parses a method declaration. Note that "func" will
// already have been consumed so we're starting from the Receiver.
// MethodDecl   = "func" Receiver MethodName ( Function | Signature ) .
// Receiver     = "(" [ identifier ] [ "*" ] BaseTypeName ")" .
func (p *Parser) parseMethodDecl() (AST, error) {
	// we already know it starts with "func"
	funcToken, _ := p.lexer.GetToken()

	// get the receiver

	return nil, NewError(p.filename, funcToken.Pos(), "unimplemented")
}

// parseFunction parses a function definition.
// Function     = Signature FunctionBody .
// FunctionBody = Block .
func (p *Parser) parseFunction() (AST, error) {
	// we already know it starts with "func"
	funcToken, _ := p.lexer.GetToken()

	return nil, NewError(p.filename, funcToken.Pos(), "unimplemented")
}

// parseSignature parses a function signature.
// Signature      = Parameters [ Result ] .
// Result         = Parameters | Type .
// Parameters     = "(" [ ParameterList [ "," ] ] ")" .
func (p *Parser) parseSignature() (AST, error) {
	// we already know it starts with "func"
	funcToken, _ := p.lexer.GetToken()

	return nil, NewError(p.filename, funcToken.Pos(), "unimplemented")
}

// parseGroupSingle parses a group of some other clause, surrounded by brackets and
// with semicolons after each entry.
func (p *Parser) parseGroupSingle(parseClause func() (AST, error), verbName string) ([]AST, error) {
	err := p.expectToken(TokenOpenBracket, "there should be a '(' here")
	if err != nil {
		return nil, err
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
		err = p.expectToken(TokenSemicolon, semiErrorMessage)
		if err != nil {
			return nil, err
		}

		asts = append(asts, newClause)
	}

	return asts, nil
}

// parseGroupMulti parses a group of some other clause, surrounded by brackets and
// with semicolons after each entry.
func (p *Parser) parseGroupMulti(parseClause func() ([]AST, error), verbName string) ([]AST, error) {
	err := p.expectToken(TokenOpenBracket, "there should be a '(' here")
	if err != nil {
		return nil, err
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
		newClauses, err := parseClause()
		if err != nil {
			return nil, err
		}

		// get a semicolon separator
		err = p.expectToken(TokenSemicolon, semiErrorMessage)
		if err != nil {
			return nil, err
		}

		asts = append(asts, newClauses...)
	}

	return asts, nil
}

// parseOptionallyQualifiedIdentifier parses an identifier with or without a package name.
// OptionallyQualifiedIdent = identifier | QualifiedIdent .
// QualifiedIdent = PackageName "." identifier .
func (p *Parser) parseOptionallyQualifiedIdentifier() (AST, error) {
	// check that it's an identifier of some sort
	tok, err := p.lexer.GetToken()
	if err != nil {
		return nil, err
	}
	if tok.TokenKind() != TokenIdentifier {
		return nil, NewError(p.filename, tok.Pos(), "if you could just put an identifier here that'd be greeeat")
	}

	ast := ASTIdentifier{tok.Pos(), "", tok.(StringToken).strVal}

	// might be followed by a '.'
	tok, err = p.lexer.PeekToken(0)
	if tok.TokenKind() == TokenDot {
		p.lexer.GetToken()

		// get a following identifier
		if tok.TokenKind() != TokenIdentifier {
			return nil, NewError(p.filename, tok.Pos(), "if you could just put an identifier here that'd be greeeat")
		}

		ast.packageName = ast.name
		ast.name = tok.(StringToken).strVal
	}

	return ast, nil
}

// expectToken parses a required token.
func (p *Parser) expectToken(tk TokenKind, message string) error {
	_, err := p.expectTokenPos(tk, message)
	return err
}

// expectTokenPos parses a required token. It returns the position of the
// token.
func (p *Parser) expectTokenPos(tk TokenKind, message string) (SrcSpan, error) {
	// get a token
	tok, err := p.lexer.GetToken()
	if err != nil {
		return tok.Pos(), err
	}
	if tok.TokenKind() != tk {
		return tok.Pos(), NewError(p.filename, tok.Pos(), message)
	}

	return tok.Pos(), nil
}
