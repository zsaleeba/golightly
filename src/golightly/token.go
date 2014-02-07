package golightly

// TokenKind indicate which type of symbol this lexical item is
type TokenKind int

const (
	// operators
	TokenAdd TokenKind = iota
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
	TokenTypeKeyword
	TokenVar

	// literals
	TokenString
	TokenRune
	TokenInt
	TokenUint
	TokenFloat32
	TokenFloat64

	// identifiers
	TokenIdentifier

	// end of source code
	TokenEndOfSource
)

// type Token is a "sum type" implemented using an interface.
// Tokens from the lexer can come with a variety of values.
// It's implemented by simpleToken, stringToken, uintToken and
// floatToken. All have the ability to have a TokenKind set,
// but each has differing ancillary values.
//
// Tokens can be created using struct initialisers.
// eg. StringToken{TokenIdentifier, "hello"}
type Token interface {
	TokenKind() TokenKind
	Pos() SrcSpan
}

type SimpleToken struct {
	pos    SrcSpan
	tt     TokenKind
}

func (st SimpleToken) TokenKind() TokenKind {
	return st.tt
}

func (st SimpleToken) Pos() SrcSpan {
	return st.pos
}

type StringToken struct {
	s SimpleToken
	strVal string
}

func (st StringToken) TokenKind() TokenKind {
	return st.s.tt
}

func (st StringToken) Pos() SrcSpan {
	return st.s.pos
}

type UintToken struct {
	s SimpleToken
	uintVal uint64
}

func (ut UintToken) TokenKind() TokenKind {
	return ut.s.tt
}

func (ut UintToken) Pos() SrcSpan {
	return ut.s.pos
}

type FloatToken struct {
	s SimpleToken
	floatVal float64
}

func (ft FloatToken) TokenKind() TokenKind {
	return ft.s.tt
}

func (st FloatToken) Pos() SrcSpan {
	return st.s.pos
}

