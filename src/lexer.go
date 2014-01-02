package golightly

import (
	"errors"
	"io"
	"unicode"
)

const (
	// operators
	TokenAdd = iota
	TokenSubtract
	TokenMultiply
	TokenDivide
	TokenModulus
	TokenBitwiseAnd
	TokenBitwiseOr
	TokenBitwiseExor
	TokenShiftLeft
	TokenShiftRight
	TokenBitClear
	TokenAddAssign
	TokenSubtractAssign
	TokenMultiplyAssign
	TokenDivideAssign
	TokenModulusAssign
	TokenBitwiseAndAssign
	TokenBitwiseOrAssign
	TokenBitwiseExorAssign
	TokenShiftLeftAssign
	TokenShiftRightAssign
	TokenBitClearAssign
	TokenLogicalAnd
	TokenLogicalOr
	TokenChannelArrow
	TokenIncrement
	TokenDecrement
	TokenEquals
	TokenLess
	TokenGreater
	TokenAssign
	TokenNot
	TokenNotEqual
	TokenLessEqual
	TokenGreaterEqual
	TokenDeclareAssign
	TokenEllipsis
	TokenOpenGroup
	TokenCloseGroup
	TokenOpenOption
	TokenCloseOption
	TokenOpenBlock
	TokenCloseBlock
	TokenComma
	TokenDot
	TokenColon
	TokenSemicolon

	// keywords
	TokenBreak
	TokenCase
	TokenChan
	TokenConst
	TokenContinue
	TokenDefault
	TokenDefer
	TokenElse
	TokenFallthrough
	TokenFor
	TokenFunc
	TokenGo
	TokenGoto
	TokenIf
	TokenImport
	TokenInterface
	TokenMap
	TokenPackage
	TokenRange
	TokenReturn
	TokenSelect
	TokenStruct
	TokenSwitch
	TokenType
	TokenVar

	// literals
	TokenString
	TokenRune
	TokenInt
	TokenUInt
	TokenFloat

	// identifiers
	TokenIdentifier
)

// a map of keywords for quick lookup
var keywords map[string]int = map[string]int{
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
	"type":        TokenType,
	"var":         TokenVar,
}

// the running state of the lexical analyser
type Lexer struct {
	sourceFile string // name of the source file
	startPos   SrcLoc // where this token started in the source
	pos        SrcLoc // where we are in the source
	lineBuf    []rune // the current source line

	tokens TokenList // the compact encoded token list
}

// LexLine lexes a line of source code and adds the tokens to the end of
// the lexed token list. The provided source should end on a line
// boundary so there are no split tokens at the end.
func (l *Lexer) LexLine(src string) error {
	// prepare for this line
	l.pos.Line++
	l.pos.Column = 0
	l.lineBuf = []rune(src)

	// get tokens until end of line
	ok := true
	for ok {
		var err error
		ok, err = l.getToken()
		if err != nil {
			return err
		}
	}

	return errors.New("unimplemented")
}

// LexReader reads all input from a Reader and lexes it until EOF.
func (l *Lexer) LexReader(r io.Reader) error {
	return errors.New("unimplemented")
}

// LexFile opens a file and lexes the entire contents.
func (l *Lexer) LexFile(filename string) error {
	return errors.New("unimplemented")
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
		if l.pos.Column >= len(l.lineBuf) {
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
			l.tokens.Add(l.startPos, token)
			return true, nil
		}

		// it must be an identifier
		l.tokens.AddString(l.startPos, TokenIdentifier, word)
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
		l.tokens.Add(l.startPos, token)
		return true, nil
	}

	// is it a string literal?
	switch ch {
	case '\'':
		l.pos.Column += 2
		err := l.getCharacterLiteral(ch2)
		return err != nil, err

	case '"', '`':
		l.pos.Column++
		err := l.getStringLiteral(ch == '`')
		return err != nil, err
	}

	return false, nil
}

// getOperator gets an operator token.
// returns the token, the number of characters absorbed and success.
func (l *Lexer) getOperator(ch, ch2 rune) (int, int, bool) {
	// operator lexing is performed as a hard-coded trie for speed.

	switch ch {
	case '+':
		switch ch2 {
		case '=':
			return TokenAddAssign, 2, true
		case '+':
			return TokenIncrement, 2, true
		default:
			return TokenAdd, 1, true
		}

	case '-':
		switch ch2 {
		case '=':
			return TokenSubtractAssign, 2, true
		case '-':
			return TokenDecrement, 2, true
		default:
			return TokenSubtract, 1, true
		}

	case '*':
		if ch2 == '=' {
			return TokenMultiplyAssign, 2, true
		} else {
			return TokenMultiply, 1, true
		}

	case '/':
		if ch2 == '=' {
			return TokenDivideAssign, 2, true
		} else {
			return TokenDivide, 1, true
		}

	case '%':
		if ch2 == '=' {
			return TokenModulusAssign, 2, true
		} else {
			return TokenModulus, 1, true
		}

	case '&':
		switch ch2 {
		case '=':
			return TokenBitwiseAndAssign, 2, true
		case '&':
			return TokenLogicalAnd, 2, true
		default:
			return TokenBitwiseAnd, 1, true
		}

	case '|':
		switch ch2 {
		case '=':
			return TokenBitwiseOrAssign, 2, true
		case '|':
			return TokenLogicalOr, 2, true
		default:
			return TokenBitwiseOr, 1, true
		}

	case '^':
		if ch2 == '=' {
			return TokenBitwiseExorAssign, 2, true
		} else {
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

			if ch3 == '=' {
				return TokenShiftLeftAssign, 3, true
			} else {
				return TokenShiftLeft, 2, true
			}
		case '=':
			return TokenLessEqual, 2, true
		case '-':
			return TokenChannelArrow, 2, true
		default:
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

			if ch3 == '=' {
				return TokenShiftRightAssign, 3, true
			} else {
				return TokenShiftRight, 2, true
			}
		case '=':
			return TokenGreaterEqual, 2, true
		default:
			return TokenGreater, 1, true
		}

	case '=':
		if ch2 == '=' {
			return TokenEquals, 2, true
		} else {
			return TokenAssign, 1, true
		}

	case '!':
		if ch2 == '=' {
			return TokenNotEqual, 2, true
		} else {
			return TokenNot, 1, true
		}

	case ':':
		if ch2 == '=' {
			return TokenDeclareAssign, 2, true
		} else {
			return TokenColon, 1, true
		}

	case '.':
		return TokenDot, 1, true
	case ',':
		return TokenComma, 1, true
	case '(':
		return TokenOpenGroup, 1, true
	case ')':
		return TokenCloseGroup, 1, true
	case '[':
		return TokenOpenOption, 1, true
	case ']':
		return TokenCloseOption, 1, true
	case '{':
		return TokenOpenBlock, 1, true
	case '}':
		return TokenCloseBlock, 1, true
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
func (l *Lexer) getNumeric() error {
	return errors.New("unimplemented")
}

// getCharacterLiteral gets a character literal.
func (l *Lexer) getCharacterLiteral(ch rune) error {
	return errors.New("unimplemented")
}

// getStringLiteral gets a string literal.
func (l *Lexer) getStringLiteral(raw bool) error {
	return errors.New("unimplemented")
}
