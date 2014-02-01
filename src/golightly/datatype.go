package golightly

// DataTypeKind indicates which type of value this is
type DataTypeKind int

const (
	// basic types
	DataTypeKindInt DataTypeKind = iota
	DataTypeKindUint
	DataTypeKindFloat
	DataTypeKindString
	DataTypeKindRune
	DataTypeKindImaginary
	DataTypeKindType

	// unary types
	DataTypeKindArray
	DataTypeKindSlice
	DataTypeKindPointer

	// struct type
	DataTypeKindStruct
)

// DataSize indicates which size value this is, in bits
type DataSize int

const (
	// operators
	DataSize16 DataSize = iota
	DataSize32
	DataSize64
	DataSizeDefault
)

// type DataType represents any Go type.
// It's a "sum type" implemented using an interface.
//
// DataType can be created using struct initialisers.
// eg. DataTypeSimple{DataTypeKindInt}
type DataType interface {
	DataTypeKind() DataTypeKind
}

// type DataTypeBasic is for "basic types" - ie. simple data types which have no sub-type
type DataTypeBasic struct {
	kind DataTypeKind
}

func (dtb *DataTypeBasic) DataTypeKind() DataTypeKind {
	return dtb.kind
}

// type DataTypeSized is for basic types which have a size - eg. int/int16/int32/int64
type DataTypeSized struct {
	kind DataTypeKind
	size DataSize
}

func (dts *DataTypeSized) DataTypeKind() DataTypeKind {
	return dts.kind
}

// type DataTypeUnary is for types which have a single sub-type
type DataTypeUnary struct {
	kind DataTypeKind
	subType *DataType
}

func (dtu *DataTypeUnary) DataTypeKind() DataTypeKind {
	return dtu.kind
}

// type DataTypeStruct is a compound data type with named fields
type DataTypeStruct struct {
	field map[string]*DataType
}

func (dtu *DataTypeStruct) DataTypeKind() DataTypeKind {
	return DataTypeKindStruct
}

