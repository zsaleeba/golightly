package golightly

// type Value is a "sum type" implemented using an interface.
// It represents literal values of any type.
//
// Values can be created using struct initialisers.
// eg. ValueString{"hello"}
type Value interface {
	isValue()
	DataType(ts *DataTypeStore) DataType
}

// type ValueInt is for signed integers
type ValueInt struct {
	typ DataType
	val int64
}

func (v ValueInt) isValue() {
}

func (v ValueInt) DataType(ts *DataTypeStore) DataType {
	return v.typ
}

// type ValueUint is for unsigned integers
type ValueUint struct {
	typ DataType
	val uint64
}

func (v ValueUint) isValue() {
}

func (v ValueUint) DataType(ts *DataTypeStore) DataType {
	return v.typ
}

// type ValueFloat is for floats
type ValueFloat struct {
	typ DataType
	val float64
}

func (v ValueFloat) isValue() {
}

func (v ValueFloat) DataType(ts *DataTypeStore) DataType {
	return v.typ
}

// type ValueRune is for runes
type ValueRune struct {
	val rune
}

func (v ValueRune) isValue() {
}

func (v ValueRune) DataType(ts *DataTypeStore) DataType {
	return ts.RuneType()
}

// type ValueString is for strings
type ValueString struct {
	val string
}

func (v ValueString) isValue() {
}

func (v ValueString) DataType(ts *DataTypeStore) DataType {
	return ts.StringType()
}

// NewValueFromToken creates a Value from a lexer Token. It assumes the
// token is a literal value type.
func NewValueFromToken(tok Token, ts *DataTypeStore) Value {
	switch tok.TokenKind() {
	case TokenKindLiteralInt:
		return ValueUint{ts.UintType(), tok.(UintToken).uintVal}
	case TokenKindLiteralFloat:
		return ValueFloat{ts.FloatType(), tok.(FloatToken).floatVal}
	case TokenKindLiteralRune:
		return ValueRune{rune(tok.(UintToken).uintVal)}
	case TokenKindLiteralString:
		return ValueString{tok.(StringToken).strVal}
	}

	return nil
}
