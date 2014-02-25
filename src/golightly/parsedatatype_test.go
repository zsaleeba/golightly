package golightly

import (
	"strings"
	"testing"
)

func setupDataTypeTest(src string) *Parser {
	lex := NewLexer()
	reader := strings.NewReader(src)
	lex.LexReader(reader, "test.go")
	ts := NewDataTypeStore()
	addSrcFile := make(chan string)
	parser := NewParser(lex, ts, addSrcFile)

	// just throw away anything we get on the addSrcFile channel.
	go func() {
		for {
			<- addSrcFile
		}
	}()

	return parser
}

func compareAST(a, b AST) bool {
	return true
}

func TestParseDataType(t *testing.T) {
	parser := setupDataTypeTest("int")
	match, ast, err := parser.parseDataType()
	if err != nil {
		t.Error("error parsing: ", err)
		return
	}
	if !match {
		t.Error("doesn't match a data type")
		return
	}
	if !compareAST(ast, ASTIdentifier{SrcSpan{SrcLoc{1,1}, SrcLoc{1,3}}, "", "int"}) {
		t.Errorf("parse failed: %s", ast)
		return
	}
}
