package golightly

// type Value is a "sum type" implemented using an interface.
// It represents literal values of any type.
//
// Values can be created using struct initialisers.
// eg. ValueString{"hello"}
type Value interface {
	isValue()
	DataType() *DataType
}

// type ValueInt is for signed integers
type ValueInt struct {
	typ *DataType
	val int64
}

func (v *ValueInt) isValue() {
}

func (v *ValueInt) DataType() *DataType {
	return v.typ
}

// type ValueUint is for unsigned integers
type ValueUint struct {
	typ *DataType
	val uint64
}

func (v *ValueUint) isValue() {
}

func (v *ValueUint) DataType() *DataType {
	return v.typ
}

// type ValueRune is for runes
type ValueRune struct {
	typ *DataType
	val rune
}

func (v *ValueRune) isValue() {
}

func (v *ValueRune) DataType() *DataType {
	return v.typ
}

// type ValueString is for strings
type ValueString struct {
	typ *DataType
	val string
}

func (v *ValueString) isValue() {
}

func (v *ValueString) DataType() *DataType {
	return v.typ
}
