package golightly

// TokenKind indicate which type of symbol this lexical item is
type TokenKind int

const (
	// operators
	TokenKindAdd TokenKind = iota
	TokenKindSubtract
	TokenKindAsterisk // can be multiplication or pointer dereference
	TokenKindDivide
	TokenKindModulus
	TokenKindBitwiseAnd
	TokenKindBitwiseOr
	TokenKindBitwiseExor
	TokenKindShiftLeft
	TokenKindShiftRight
	TokenKindBitClear
	TokenKindAddAssign
	TokenKindSubtractAssign
	TokenKindMultiplyAssign
	TokenKindDivideAssign
	TokenKindModulusAssign
	TokenKindBitwiseAndAssign
	TokenKindBitwiseOrAssign
	TokenKindBitwiseExorAssign
	TokenKindShiftLeftAssign
	TokenKindShiftRightAssign
	TokenKindBitClearAssign
	TokenKindLogicalAnd
	TokenKindLogicalOr
	TokenKindChannelArrow
	TokenKindIncrement
	TokenKindDecrement
	TokenKindEquals
	TokenKindLess
	TokenKindGreater
	TokenKindAssign
	TokenKindNot
	TokenKindNotEqual
	TokenKindLessEqual
	TokenKindGreaterEqual
	TokenKindDeclareAssign
	TokenKindEllipsis
	TokenKindOpenBracket
	TokenKindCloseBracket
	TokenKindOpenSquareBracket
	TokenKindCloseSquareBracket
	TokenKindOpenBrace
	TokenKindCloseBrace
	TokenKindComma
	TokenKindDot
	TokenKindColon
	TokenKindSemicolon

	// keywords
	TokenKindBreak
	TokenKindCase
	TokenKindChan
	TokenKindConst
	TokenKindContinue
	TokenKindDefault
	TokenKindDefer
	TokenKindElse
	TokenKindFallthrough
	TokenKindFor
	TokenKindFunc
	TokenKindGo
	TokenKindGoto
	TokenKindIf
	TokenKindImport
	TokenKindInterface
	TokenKindMap
	TokenKindPackage
	TokenKindRange
	TokenKindReturn
	TokenKindSelect
	TokenKindStruct
	TokenKindSwitch
	TokenKindTypeKeyword
	TokenKindVar

	// data type keywords
	TokenKindBool
	TokenKindUint
	TokenKindUint8
	TokenKindUint16
	TokenKindUint32
	TokenKindUint64
	TokenKindUintPtr
	TokenKindInt
	TokenKindInt8
	TokenKindInt16
	TokenKindInt32
	TokenKindInt64
	TokenKindFloat32
	TokenKindFloat64
	TokenKindComplex64
	TokenKindComplex128
	TokenKindByte
	TokenKindRune
	TokenKindString
	TokenKindError

	// identifiers
	TokenKindIdentifier

	// literals
	TokenKindLiteralInt
	TokenKindLiteralFloat
	TokenKindLiteralRune
	TokenKindLiteralString

	// end of source code
	TokenKindEndOfSource
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
	pos SrcSpan
	tt  TokenKind
}

func (st SimpleToken) TokenKind() TokenKind {
	return st.tt
}

func (st SimpleToken) Pos() SrcSpan {
	return st.pos
}

type StringToken struct {
	s      SimpleToken
	strVal string
}

func (st StringToken) TokenKind() TokenKind {
	return st.s.tt
}

func (st StringToken) Pos() SrcSpan {
	return st.s.pos
}

type UintToken struct {
	s       SimpleToken
	uintVal uint64
}

func (ut UintToken) TokenKind() TokenKind {
	return ut.s.tt
}

func (ut UintToken) Pos() SrcSpan {
	return ut.s.pos
}

type FloatToken struct {
	s        SimpleToken
	floatVal float64
}

func (ft FloatToken) TokenKind() TokenKind {
	return ft.s.tt
}

func (st FloatToken) Pos() SrcSpan {
	return st.s.pos
}
