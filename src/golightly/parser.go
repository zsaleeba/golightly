package golightly

// type Parser controls parsing of a token stream into an AST
type Parser struct {
	filename    string // the name of the file being parsed
	tokenLoc    SrcLoc // the location in the file of the current token
	tokenEndLoc SrcLoc // the location in the file of the end of the current token

	tokenQueue  []Token // a buffer of the incoming tokens

	out chan *AST // the stream of ASTs is sent out through this channel
}

func NewParser() *Parser {
	p := new(Parser)
	p.out = make(chan *AST)

	return p
}

// Parse runs the parser and breaks the program down into an Abstract Syntax Tree
func (p *Parser) Parse() error {
	return nil
}
