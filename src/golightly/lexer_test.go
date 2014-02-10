package golightly

import (
	"testing"
)

func TestLexerLexLine(t *testing.T) {
	l := NewLexer()
	l.Init("-")
	addLexLine(t, l, "package golightly")
	addLexLine(t, l, "")
	addLexLine(t, l, "import (")
	addLexLine(t, l, "    \"testing\"")
	addLexLine(t, l, ")")
	addLexLine(t, l, "")
	addLexLine(t, l, "i := 42;")
	addLexLine(t, l, "j := 7.2;")
	addLexLine(t, l, "k += 'X';")

	// now try to get them back out
	tl := l.tokens
	tl.StartReading()

	err := checkToken(tl, 1, 1, TokenPackage)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 1, 9, TokenIdentifier, "golightly")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 3, 1, TokenImport)
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 3, 8, TokenOpenGroup)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 4, 5, TokenString, "testing")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 5, 1, TokenCloseGroup)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 7, 1, TokenIdentifier, "i")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 7, 3, TokenDeclareAssign)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenUint(tl, 7, 6, TokenLiteralInt, 42)
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 7, 8, TokenSemicolon)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 8, 1, TokenIdentifier, "j")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 8, 3, TokenDeclareAssign)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenFloat(tl, 8, 6, TokenLiteralFloat, 7.2)
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 8, 9, TokenSemicolon)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 9, 1, TokenIdentifier, "k")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 9, 3, TokenAddAssign)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenUint(tl, 9, 6, TokenLiteralRune, uint64('X'))
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 9, 9, TokenSemicolon)
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 10, 1, TokenEndOfSource)
	if err != nil {
		t.Error(err)
	}
}

func addLexLine(t *testing.T, l *Lexer, src string) {
	err := l.LexLine(src)
	if err != nil {
		t.Errorf("LexLine() failed: '%s', %s", src, err)
	}
}

func TestLexerGetWord(t *testing.T) {
	l := setupLexerTest("hello")
	if l.getWord() != "hello" {
		t.Error("getWord() failed")
	}

	l = setupLexerTest("hello ")
	if l.getWord() != "hello" {
		t.Error("getWord() failed")
	}

	l = setupLexerTest("hello<")
	if l.getWord() != "hello" {
		t.Error("getWord() failed")
	}
}

func TestLexerGetNumericInteger(t *testing.T) {
	// integer with no trailing character
	l := setupLexerTest("12345")
	err := l.getNumeric()
	if err != nil {
		t.Errorf("getNumeric() failed: %s", err)
	}

	l.tokens.StartReading()
	tok, _ := l.tokens.GetToken()
	if tok != TokenUint || l.tokens.GetValueUint64() != 12345 {
		t.Error("getNumeric() failed")
	}

	// now with a trailing character
	l = setupLexerTest("36212;")
	err = l.getNumeric()
	if err != nil {
		t.Errorf("getNumeric() failed: %s", err)
	}

	l.tokens.StartReading()
	tok, _ = l.tokens.GetToken()
	if tok != TokenUint || l.tokens.GetValueUint64() != 36212 {
		t.Error("getNumeric() failed")
	}
}

func TestLexerGetNumericFloat(t *testing.T) {
	// integer with no trailing character
	l := setupLexerTest("12.345")
	err := l.getNumeric()
	if err != nil {
		t.Errorf("getNumeric() failed: %s", err)
	}

	l.tokens.StartReading()
	tok, _ := l.tokens.GetToken()
	if tok != TokenFloat64 || l.tokens.GetValueFloat64() != 12.345 {
		t.Error("getNumeric() failed")
	}

	// now with a trailing character
	l = setupLexerTest("1.469e1;")
	err = l.getNumeric()
	if err != nil {
		t.Errorf("getNumeric() failed: %s", err)
	}

	l.tokens.StartReading()
	tok, _ = l.tokens.GetToken()
	if tok != TokenFloat64 || l.tokens.GetValueFloat64() != 1.469e1 {
		t.Error("getNumeric() failed")
	}
}

func TestLexerGetRuneLiteral(t *testing.T) {
	l := setupLexerTest("'a'")
	err := l.getRuneLiteral()
	if err != nil {
		t.Errorf("getRuneLiteral() failed: %s", err)
	}

	l.tokens.StartReading()
	tok, _ := l.tokens.GetToken()
	if tok != TokenRune || l.tokens.GetValueUint64() != uint64('a') {
		t.Error("getRuneLiteral() failed")
	}
}

func TestLexerGetStringLiteral(t *testing.T) {
	l := setupLexerTest("\"hello\"")
	err := l.getStringLiteral()
	if err != nil {
		t.Errorf("getStringLiteral() failed: %s", err)
	}

	l.tokens.StartReading()
	tok, _ := l.tokens.GetToken()
	if tok != TokenString || l.tokens.GetValueString() != "hello" {
		t.Error("getStringLiteral() failed")
	}

	l = setupLexerTest("`hello`")
	err = l.getStringLiteral()
	if err != nil {
		t.Errorf("getStringLiteral() failed: %s", err)
	}

	l.tokens.StartReading()
	tok, _ = l.tokens.GetToken()
	if tok != TokenString || l.tokens.GetValueString() != "hello" {
		t.Error("getStringLiteral() failed")
	}
}

func setupLexerTest(source string) *Lexer {
	l := NewLexer()
	l.Init("-")
	l.lineBuf = []rune(" " + source)

	return l
}
