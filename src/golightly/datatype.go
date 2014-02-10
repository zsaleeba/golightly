package golightly

import "sync"

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

// DataSize indicates which size value this is.
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

// type DataTypeBasic is for "basic types" - ie. simple data types which have no sub-type.
type DataTypeBasic struct {
	kind DataTypeKind
}

func (dtb DataTypeBasic) DataTypeKind() DataTypeKind {
	return dtb.kind
}

// type DataTypeSized is for basic types which have a size - eg. int/int16/int32/int64.
type DataTypeSized struct {
	kind DataTypeKind
	size DataSize
}

func (dts DataTypeSized) DataTypeKind() DataTypeKind {
	return dts.kind
}

// type DataTypeUnary is for types which have a single sub-type.
type DataTypeUnary struct {
	kind    DataTypeKind
	subType *DataType
}

func (dtu DataTypeUnary) DataTypeKind() DataTypeKind {
	return dtu.kind
}

// type DataTypeStruct is a compound data type with named fields.
type DataTypeStruct struct {
	field map[string]*DataType
}

func (dtu DataTypeStruct) DataTypeKind() DataTypeKind {
	return DataTypeKindStruct
}

// type DataTypeStore is a store of all the data types in the system. Each
// unique data type will be stored only once and a reference to it always
// returns the same pointer so pointer comparison can be used on types.
type DataTypeStore struct {
	// a map of type names to types
	nameMap      map[string]DataType
	nameMapMutex sync.RWMutex

	// standard types
	intType    DataType
	uintType   DataType
	floatType  DataType
	runeType   DataType
	stringType DataType
}

// NewDataTypeStore creates a new data type store.
func NewDataTypeStore() *DataTypeStore {
	ts := new(DataTypeStore)

	// add the predefined data types
	ts.intType = DataTypeSized{DataTypeKindInt, DataSizeDefault}
	ts.uintType = DataTypeSized{DataTypeKindUint, DataSizeDefault}
	ts.floatType = DataTypeSized{DataTypeKindFloat, DataSizeDefault}
	ts.runeType = DataTypeBasic{DataTypeKindRune}
	ts.stringType = DataTypeBasic{DataTypeKindString}

	ts.nameMapMutex.Lock()
	ts.nameMap["int"] = ts.intType
	ts.nameMap["uint"] = ts.uintType
	ts.nameMap["float"] = ts.floatType
	ts.nameMap["rune"] = ts.runeType
	ts.nameMap["string"] = ts.stringType
	ts.nameMapMutex.Unlock()

	return ts
}

// methods to get all the predefined types.
func (ts *DataTypeStore) IntType() DataType {
	return ts.intType
}
func (ts *DataTypeStore) UintType() DataType {
	return ts.uintType
}
func (ts *DataTypeStore) FloatType() DataType {
	return ts.floatType
}
func (ts *DataTypeStore) RuneType() DataType {
	return ts.runeType
}
func (ts *DataTypeStore) StringType() DataType {
	return ts.stringType
}
