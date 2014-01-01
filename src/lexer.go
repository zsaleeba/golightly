package golightly

import (
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
	TokenNotEquals
	TokenLessEquals
	TokenGreaterEquals
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
)

// a map of keywords for quick lookup
var keywords map[string]int = {
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
struct Lexer {
	sourceFile string // name of the source file
	pos SrcLoc        // where we are in the source
	lineBuf []rune    // the current source line
	
	tokens TokenList  // the compact encoded token list
}

// LexLine lexes a line of source code and adds the tokens to the end of
// the lexed token list. The provided source should end on a line
// boundary so there are no split tokens at the end. 
func (l *Lexer) LexLine(src string) error {
	// prepare for this line
	l.line++
	l.pos.column = 0
	l.lineBuf = src

	// get tokens until end of line	
	ok := true
	for ok {
		token, ok, err := getToken()
		if err != nil {
			return err
		}
		
		tokens.Add(token)
	}
	
	return errors.New("unimplemented")
}

// LexReader reads all input from a Reader and lexes it until EOF.
func (l *Lexer) LexReader(r Reader) error {
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
	if l.lineOffset >= len(l.lineBuf) {
		return false, nil
	}
	
	// skip whitespace
	ch := l.lineBuf[l.lineOffset]
	for ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
		l.lineOffset++
		if l.lineOffset >= len(l.lineBuf) {
			return false, nil    // end of line
		}
		ch = l.lineBuf[l.lineOffset]
	}
	
	// is it an identifier?
	if unicode.IsLetter(ch) || ch == '_' {
		// get the word
		word := getWord()
		
		// is it a keyword?
		token, ok := keywords[word]
		if ok {
			l.tokens.Add(token)
			return true, nil
		}
		
		// it must be an identifier
		l.tokens.AddString(TokenIdentifier, word)
		return true, nil
	}
	
	// is it a numeric literal?
	var ch2 rune
	if l.lineOffset+1 < len(l.lineBuf) {
		ch2 = l.lineBuf[l.lineOffset+1]
	}
	
	if unicode.IsDigit(ch) || (ch == '.' && unicode.IsDigit(ch2)) {
		err := getNumeric()
		return true, err
	} 
	
	// is it an operator?
	token, runes, ok := getOperator(ch, ch2)
	if ok {
		l.pos.column += runes
		l.tokens.Add(token)
		return true, nil
	}
	
	// is it a string literal?
	switch ch {
		case '\'':
			l.pos.column += 2
			err := getCharacterLiteral(ch2)
			return err != nil, err
		
		case '"':
		case '`':
			l.pos.column++
			token, word, err := getStringLiteral(ch == '`')
			
			return token, true, err
	}
	
	return 0, false
}

// getOperator gets an operator token.
// returns the token, the number of characters absorbed and success.
func (l *Lexer) getOperator(ch, ch2 rune) (int, int, bool) {
	// operator lexing is performed as a hard-coded trie for speed.

	switch ch {
		case '+':
			switch ch2 {
				case '=': return TokenAddAssign, 2, true
				case '+': return TokenIncrement, 2, true
				default:  return TokenAdd, 1, true
			}
			
		case '-':
			switch ch2 {
				case '=': return TokenSubtractAssign, 2, true
				case '-': return TokenDecrement, 2, true
				default:  return TokenSubtract, 1, true
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
				case '=': return TokenBitwiseAndAssign, 2, true
				case '&': return TokenLogicalAnd, 2, true
				default:  return TokenBitwiseAnd, 1, true
			}
			
		case '|':
			switch ch2 {
				case '=': return TokenBitwiseOrAssign, 2, true
				case '|': return TokenLogicalOr, 2, true
				default:  return TokenBitwiseOr, 1, true
			}
			
		case '^':
			if ch2 == '=' {
				return TokenExorAssign, 2, true
			} else {
				return TokenExor, 1, true
			}
			
		case '<':
			switch ch2 {
				// look ahead another character
				var ch3 rune
				if l.pos.column+2 < len(l.lineBuf) {
					ch3 = l.lineBuf[l.pos.column+2]
				}

				case '<': 
					if ch3 == '=' {
						return TokenShiftLeftAssign, 3, true
					} else {
						return TokenShiftLeft, 2, true
					}
				case '=': return TokenLessEqual, 2, true
				case '-': return TokenChannelArrow, 2, true
				default:  return TokenLess, 1, true
			}
			
		case '>':
			switch ch2 {
				// look ahead another character
				var ch3 rune
				if l.pos.column+2 < len(l.lineBuf) {
					ch3 = l.lineBuf[l.pos.column+2]
				}

				case '>': 
					if ch3 == '=' {
						return TokenShiftRightAssign, 3, true
					} else {
						return TokenShiftRight, 2, true
					}
				case '=': return TokenGreaterEqual, 2, true
				default:  return TokenGreater, 1, true
			}
			
		case '=':
			if ch2 == '=' {
				return TokenEquals, 2, true
			} else {
				return TokenAssign, 1, true
			}
			
		case '!':
			if ch2 == '=' {
				return TokenNotAssign, 2, true
			} else {
				return TokenNot, 1, true
			}
			
		case ':':
			if ch2 == '=' {
				return TokenDeclareAssign, 2, true
			} else {
				return TokenColon, 1, true
			}
			
		case '.': return TokenDot, 1, true
		case ',': return TokenComma, 1, true
		case '(': return TokenOpenGroup, 1, true
		case ')': return TokenCloseGroup, 1, true
		case '[': return TokenOpenOption, 1, true
		case ']': return TokenCloseOption, 1, true
		case '{': return TokenOpenBlock, 1, true
		case '}': return TokenCloseBlock, 1, true
	}
	
	return 0, 0, false
}
