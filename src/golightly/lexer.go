package golightly

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

// a map of keywords for quick lookup
var keywords map[string]TokenType = map[string]TokenType{
	"break":       TokenBreak,
	"case":        TokenCase,
	"chan":        TokenChan,
	"const":       TokenConst,
	"continue":    TokenContinue,
	"default":     TokenDefault,
	"defer":       TokenDefer,
	"else":        TokenElse,
	"fallthrough": TokenFallthrough,
	"for":         TokenFor,
	"func":        TokenFunc,
	"go":          TokenGo,
	"goto":        TokenGoto,
	"if":          TokenIf,
	"import":      TokenImport,
	"interface":   TokenInterface,
	"map":         TokenMap,
	"package":     TokenPackage,
	"range":       TokenRange,
	"return":      TokenReturn,
	"select":      TokenSelect,
	"struct":      TokenStruct,
	"switch":      TokenSwitch,
	"type":        TokenTypeKeyword,
	"var":         TokenVar,
}

// the running state of the lexical analyser
type Lexer struct {
	sourceFile string // name of the source file
	startPos   SrcLoc // where this token started in the source
	pos        SrcLoc // where we are in the source
	lineBuf    []rune // the current source line

	out chan Token // the token stream is sent out through this channel
}

// the buffer size of the lexer output channel
const lexerTokenChannelBuffers = 5

// NewLexer creates a new lexer object
func NewLexer() *Lexer {
	l := new(Lexer)
	l.out = make(chan Token, lexerTokenChannelBuffers)
	l.Init("-")
	return l
}

// Init initialises the lexer before using LexLine.
func (l *Lexer) Init(filename string) {
	l.pos.Line = 1
	l.pos.Column = 1
	l.startPos = l.pos
	l.sourceFile = filename
}

// Tokens returns a channel of tokens as output from the lexer.
func (l *Lexer) Tokens() chan Token {
	return l.out
}

// LexLine lexes a line of source code and adds the tokens to the end of
// the lexed token list. The provided source should end on a line
// boundary so there are no split tokens at the end.
func (l *Lexer) LexLine(src string) error {
	// prepare for this line
	defer func() {
		l.pos.Column = 1
		l.pos.Line++
	}()

	// since columns are 1-based we add a spurious character at the start so lineBuf[1] is the first character
	l.lineBuf = []rune(" " + src)

	// get tokens until end of line
	ok := true
	for ok {
		var err error
		ok, err = l.getToken()
		if err != nil {
			return err
		}
	}

	return nil
}

// LexReader reads all input from a Reader and lexes it until EOF.
func (l *Lexer) LexReader(r io.Reader, filename string) error {
	// start afresh
	l.Init(filename)

	// get lines until EOF
	scanner := bufio.NewScanner(r)
	var err error
	for scanner.Scan() {
		// get the line
		l.lineBuf = []rune(" " + scanner.Text())

		// tokenise the line
		var ok bool
		for ok, err = l.getToken(); ok && err == nil; {
		}

		if err != nil {
			return err
		}
	}

	// check for any line scanner errors
	err = scanner.Err()
	return err
}

// LexFile opens a file and lexes the entire contents.
func (l *Lexer) LexFile(filename string) error {
	// open the file
	inFile, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	// now lex it
	return l.LexReader(reader, filename)
}

// getToken gets the next token from the line buffer.
// adds the token to the token list.
// returns success and an error. success is false at end of line.
func (l *Lexer) getToken() (bool, error) {
	// are there any characters left?
	if l.pos.Column >= len(l.lineBuf) {
		return false, nil
	}

	// skip leading whitespace
	ch := l.lineBuf[l.pos.Column]
	for ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
		l.pos.Column++
		if l.pos.Column > len(l.lineBuf) {
			return false, nil // end of line
		}
		ch = l.lineBuf[l.pos.Column]
	}

	l.startPos = l.pos

	// is it an identifier?
	if unicode.IsLetter(ch) || ch == '_' {
		// get the word
		word := l.getWord()

		// is it a keyword?
		token, ok := keywords[word]
		if ok {
			l.out <- SimpleToken{l.startPos, token}
			return true, nil
		}

		// it must be an identifier
		l.out <- StringToken{l.startPos, TokenIdentifier, word}
		return true, nil
	}

	// is it a numeric literal?
	var ch2 rune
	if l.pos.Column+1 < len(l.lineBuf) {
		ch2 = l.lineBuf[l.pos.Column+1]
	}

	if unicode.IsDigit(ch) || (ch == '.' && unicode.IsDigit(ch2)) {
		err := l.getNumeric()
		return true, err
	}

	// is it an operator?
	token, runes, ok := l.getOperator(ch, ch2)
	if ok {
		l.pos.Column += runes
		l.out <- SimpleToken{l.startPos, token}
		return true, nil
	}

	// is it a string literal?
	switch ch {
	case '\'':
		err := l.getRuneLiteral()
		return err == nil, err

	case '"', '`':
		err := l.getStringLiteral()
		return err == nil, err
	}

	return false, errors.New(fmt.Sprintf("illegal character '%c' (0x%02x)", ch, ch))
}

// getOperator gets an operator token.
// returns the token, the number of characters absorbed and success.
func (l *Lexer) getOperator(ch, ch2 rune) (TokenType, int, bool) {
	// operator lexing is performed as a hard-coded trie for speed.

	switch ch {
	case '+':
		switch ch2 {
		case '=': // '+='
			return TokenAddAssign, 2, true
		case '+': // '++'
			return TokenIncrement, 2, true
		default: // '+'
			return TokenAdd, 1, true
		}

	case '-':
		switch ch2 {
		case '=': // '-='
			return TokenSubtractAssign, 2, true
		case '-': // '--'
			return TokenDecrement, 2, true
		default: // '-'
			return TokenSubtract, 1, true
		}

	case '*':
		if ch2 == '=' { // '*='
			return TokenMultiplyAssign, 2, true
		} else { // '*'
			return TokenMultiply, 1, true
		}

	case '/':
		if ch2 == '=' { // '/='
			return TokenDivideAssign, 2, true
		} else { // '/'
			return TokenDivide, 1, true
		}

	case '%':
		if ch2 == '=' { // '%='
			return TokenModulusAssign, 2, true
		} else { // '%'
			return TokenModulus, 1, true
		}

	case '&':
		switch ch2 {
		case '=': // '&='
			return TokenBitwiseAndAssign, 2, true
		case '&': // '&&'
			return TokenLogicalAnd, 2, true
		default: // '&'
			return TokenBitwiseAnd, 1, true
		}

	case '|':
		switch ch2 {
		case '=': // '|='
			return TokenBitwiseOrAssign, 2, true
		case '|': // '||'
			return TokenLogicalOr, 2, true
		default: // '|'
			return TokenBitwiseOr, 1, true
		}

	case '^':
		if ch2 == '=' { // '^='
			return TokenBitwiseExorAssign, 2, true
		} else { // '^'
			return TokenBitwiseExor, 1, true
		}

	case '<':
		switch ch2 {
		case '<':
			// look ahead another character
			var ch3 rune
			if l.pos.Column+2 < len(l.lineBuf) {
				ch3 = l.lineBuf[l.pos.Column+2]
			}

			if ch3 == '=' { // '<<='
				return TokenShiftLeftAssign, 3, true
			} else { // '<<'
				return TokenShiftLeft, 2, true
			}
		case '=': // '<='
			return TokenLessEqual, 2, true
		case '-': // '<-'
			return TokenChannelArrow, 2, true
		default: // '<'
			return TokenLess, 1, true
		}

	case '>':
		switch ch2 {
		case '>':
			// look ahead another character
			var ch3 rune
			if l.pos.Column+2 < len(l.lineBuf) {
				ch3 = l.lineBuf[l.pos.Column+2]
			}

			if ch3 == '=' { // '>>='
				return TokenShiftRightAssign, 3, true
			} else { // '>>'
				return TokenShiftRight, 2, true
			}
		case '=': // '>='
			return TokenGreaterEqual, 2, true
		default: // '>'
			return TokenGreater, 1, true
		}

	case '=':
		if ch2 == '=' { // '=='
			return TokenEquals, 2, true
		} else { // '='
			return TokenAssign, 1, true
		}

	case '!':
		if ch2 == '=' { // '!='
			return TokenNotEqual, 2, true
		} else { // '!'
			return TokenNot, 1, true
		}

	case ':':
		if ch2 == '=' { // ':='
			return TokenDeclareAssign, 2, true
		} else { // ':'
			return TokenColon, 1, true
		}

	case '.': // '.'
		return TokenDot, 1, true
	case ',': // ','
		return TokenComma, 1, true
	case '(': // '('
		return TokenOpenGroup, 1, true
	case ')': // ')'
		return TokenCloseGroup, 1, true
	case '[': // '['
		return TokenOpenOption, 1, true
	case ']': // ']'
		return TokenCloseOption, 1, true
	case '{': // '{'
		return TokenOpenBlock, 1, true
	case '}': // '}'
		return TokenCloseBlock, 1, true
	case ';': // ';'
		return TokenSemicolon, 1, true
	}

	return 0, 0, false
}

// getWord gets an identifier. returns the word.
func (l *Lexer) getWord() string {
	// get character until end of line
	for ; l.pos.Column < len(l.lineBuf); l.pos.Column++ {
		ch := l.lineBuf[l.pos.Column]

		// done at end of word
		if !unicode.IsLetter(ch) && ch != '_' {
			return string(l.lineBuf[l.startPos.Column:l.pos.Column])
		}
	}

	// reached end of line
	return string(l.lineBuf[l.startPos.Column:l.pos.Column])
}

// getNumeric gets a number.
// XXX - this is currently a quickie version. This should be reimplemented fully according to spec later.
func (l *Lexer) getNumeric() error {
	// scan for a non-digit character
	var col int
	for col = l.pos.Column; col < len(l.lineBuf) && unicode.IsDigit(l.lineBuf[col]); col++ {
	}

	// is the next character a "." or "e"? If so, it's a float.
	if col < len(l.lineBuf) && (l.lineBuf[col] == '.' || l.lineBuf[col] == 'e') {
		// it's a float, scan for the end
		for col = l.pos.Column; col < len(l.lineBuf) && (unicode.IsDigit(l.lineBuf[col]) || l.lineBuf[col] == '.' || l.lineBuf[col] == 'e'); col++ {
		}

		// parse the float
		v, err := strconv.ParseFloat(string(l.lineBuf[l.pos.Column:col]), 128)
		if err != nil {
			return err
		}

		l.out <- FloatToken{l.startPos, TokenFloat64, v}
		l.pos.Column = col
		return nil
	} else {
		// it's an int, parse it
		v, err := strconv.ParseUint(string(l.lineBuf[l.pos.Column:col]), 10, 64)
		if err != nil {
			return err
		}

		l.out <- UintToken{l.startPos, TokenUint, v}
		l.pos.Column = col
		return nil
	}
}

// getRuneLiteral gets a single character rune literal.
// XXX - this is currently a quickie version. This should be reimplemented fully according to spec later.
func (l *Lexer) getRuneLiteral() error {
	if l.pos.Column+2 >= len(l.lineBuf) {
		return errors.New("incomplete rune literal")
	}
	ch := l.lineBuf[l.pos.Column+1]
	l.out <- UintToken{l.startPos, TokenRune, uint64(ch)}
	if l.lineBuf[l.pos.Column+2] != '\'' {
		return errors.New("expected closing single quote in rune literal")
	}

	l.pos.Column += 3

	return nil
}

// getStringLiteral gets a string literal.
// XXX - this is currently a quickie version. This should be reimplemented fully according to spec later.
func (l *Lexer) getStringLiteral() error {
	quote := l.lineBuf[l.pos.Column]
	sCol := l.pos.Column + 1
	var col int

	for col = sCol; col < len(l.lineBuf) && l.lineBuf[col] != quote; col++ {
	}

	if col == len(l.lineBuf) {
		return errors.New("no closing quote")
	}

	l.out <- StringToken{l.startPos, TokenString, string(l.lineBuf[sCol:col])}
	l.pos.Column = col + 1
	return nil
}
