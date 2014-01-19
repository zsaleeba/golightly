package golightly

import (
	"bytes"
	"encoding/binary"
	"io"
)

// an encoded, compact form of the tokens
type TokenList struct {
	last    SrcLoc
	tokens  *bytes.Buffer
	symbols []string
	reader  io.Reader

	// temporary storage used while tokenising the program
	symMap map[string]int

	// literal values associated with a token
	strVal   string
	intVal   int64
	uintVal  uint64
	floatVal float64
}

const TokenFlagInt16 = 0xfd
const TokenFlagInt32 = 0xfe
const TokenFlagInt64 = 0xff

const tokenListInitialSymbols = 32

var endian = binary.LittleEndian

func NewTokenList(filename string) *TokenList {
	tl := new(TokenList)
	tl.last.Line = 1
	tl.last.Column = 1
	tl.tokens = new(bytes.Buffer)
	tl.symbols = make([]string, 0, tokenListInitialSymbols)

	tl.symMap = make(map[string]int)

	return tl
}

func (tl *TokenList) Add(pos SrcLoc, token Token) {
	tl.EncodeLoc(pos)
	tl.tokens.WriteByte(byte(token))
}

func (tl *TokenList) AddInt(pos SrcLoc, token Token, val int64) {
	tl.EncodeLoc(pos)
	tl.tokens.WriteByte(byte(token))
	tl.EncodeInt64(val)
}

func (tl *TokenList) AddUInt(pos SrcLoc, token Token, val uint64) {
	tl.EncodeLoc(pos)
	tl.tokens.WriteByte(byte(token))
	tl.EncodeUint64(val)
}

func (tl *TokenList) AddString(pos SrcLoc, token Token, str string) {
	tl.EncodeLoc(pos)
	tl.tokens.WriteByte(byte(token))

	// put the symbol in the symbol slice
	offset, ok := tl.symMap[str]
	if !ok {
		offset = len(tl.symMap)
		tl.symMap[str] = offset
		tl.symbols = append(tl.symbols, str)
	}

	tl.EncodeUint64(uint64(offset))
}

func (tl *TokenList) AddFloat(pos SrcLoc, val float64) {
	tl.EncodeLoc(pos)
	v32 := float32(val)
	if float64(v32) == val {
		// can be represented as a float32
		tl.tokens.WriteByte(byte(TokenFloat32))
		binary.Write(tl.tokens, endian, v32)
	} else {
		// we need the full 64 bits
		tl.tokens.WriteByte(byte(TokenFloat64))
		binary.Write(tl.tokens, endian, val)
	}
}

// EncodeLoc stores the location of this token as a delta from the
// previous token. If this token is on the same line as the previous
// token the number of columns to advance is stored as a positive
// integer. If the line has changed the negative of the absolute
// column is stored and then the number of lines to advance is stored.
func (tl *TokenList) EncodeLoc(pos SrcLoc) {
	if tl.last.Line == pos.Line {
		// it's on the same line so just output a delta for the column
		tl.EncodeInt64(int64(pos.Column - tl.last.Column)) // positive
		tl.last.Column = pos.Column
	} else {
		// new line so encode a delta for the line and an absolute column
		tl.EncodeInt64(-int64(pos.Column)) // negative
		tl.EncodeUint64(uint64(pos.Line - tl.last.Line))
		tl.last = pos
	}
}

// DecodeLoc decodes a SrcLoc encoded as described above.
func (tl *TokenList) DecodeLoc() SrcLoc {
	v1 := tl.DecodeInt64()
	if v1 >= 0 {
		tl.last.Column += int(v1)
	} else {
		tl.last.Column = -int(v1)
		lineInc := tl.DecodeUint64()
		tl.last.Line += int(lineInc)
	}

	return tl.last
}

// EncodeInt64 encodes a signed number using a variable precision method.
// The LSB being set indicates a negative number.
//  positive values are shifted left one bit and encoded using EncodeUint64.
//  negative values are flagged as negative, negated and encoded using
//     EncodeUint64. This ensures that small negative values don't use much
//     space.
func (tl *TokenList) EncodeInt64(val int64) {
	if val < 0 {
		// negative number
		tl.EncodeUint64((uint64(val^0x7fffffffffffffff) << 1) | 0x01)
	} else {
		// positive number
		tl.EncodeUint64(uint64(val << 1))
	}
}

// DecodeInt64 decodes a number encoded as described above.
func (tl *TokenList) DecodeInt64() int64 {
	val := tl.DecodeUint64()
	if val&0x01 != 0 {
		// negative
		return int64((val >> 1) ^ 0xffffffffffffffff)
	} else {
		// positive
		return int64(val >> 1)
	}
}

// EncodeUint64 encodes a number using a variable precision method.
//  values < 0xfd are simply stored as a byte.
//  value >= 0xfd and < 0x10000 are stored as flag 0xfd and a 16 bit little
//     endian value.
//  value >= 0x10000 and < 0x100000000 are stored as flag 0xfe and a 32 bit
//     little endian value.
//  value >= 0x100000000 are stored as flag 0xff and a 64 bit little endian
//     value.
func (tl *TokenList) EncodeUint64(val uint64) {
	if val < TokenFlagInt16 {
		binary.Write(tl.tokens, endian, byte(val))
	} else if val < 0x10000 {
		binary.Write(tl.tokens, endian, byte(TokenFlagInt16))
		binary.Write(tl.tokens, endian, uint16(val))
	} else if val < 0x100000000 {
		binary.Write(tl.tokens, endian, byte(TokenFlagInt32))
		binary.Write(tl.tokens, endian, uint32(val))
	} else {
		binary.Write(tl.tokens, endian, byte(TokenFlagInt64))
		binary.Write(tl.tokens, endian, val)
	}
}

// DecodeUint64 decodes a number encoded as described above.
func (tl *TokenList) DecodeUint64() uint64 {
	var b byte
	err := binary.Read(tl.reader, endian, &b)
	if err != nil {
		return 0
	}

	if b < TokenFlagInt16 {
		return uint64(b)
	}

	var result uint64
	switch b {
	case TokenFlagInt16:
		var val uint16
		err = binary.Read(tl.reader, endian, &val)
		result = uint64(val)

	case TokenFlagInt32:
		var val uint32
		err = binary.Read(tl.reader, endian, &val)
		result = uint64(val)

	case TokenFlagInt64:
		err = binary.Read(tl.reader, endian, &result)
	}

	if err != nil {
		return 0
	}

	return result
}

// StartReading resets the read position to the start of the TokenList.
// This should be called before using GetToken, and called again to re-read
// the tokens.
func (tl *TokenList) StartReading() {
	tl.reader = bytes.NewReader(tl.tokens.Bytes())
	tl.last.Line = 1
	tl.last.Column = 1
}

// GetToken gets a single token from the TokenList. It returns
// TokenEndOfSource at the end. Some tokens have associated literal values
// which can be retrieved using GetValueXXX. The location of the token in
// the source code is set in <loc>.
func (tl *TokenList) GetToken() (Token, SrcLoc) {
	// get the source location
	loc := tl.DecodeLoc()

	// get the token byte
	var b byte
	err := binary.Read(tl.reader, endian, &b)
	if err != nil {
		return TokenEndOfSource, loc
	}

	token := Token(b)
	if b < byte(TokenString) {
		// keywords and operators have no value
		return token, loc
	} else {
		switch token {
		// literals
		case TokenString, TokenIdentifier:
			symIndex := int(tl.DecodeUint64())
			if symIndex >= len(tl.symbols) {
				return TokenEndOfSource, loc
			}
			tl.strVal = tl.symbols[symIndex]

		case TokenRune, TokenUint:
			tl.uintVal = tl.DecodeUint64()

		case TokenInt:
			tl.intVal = tl.DecodeInt64()

		case TokenFloat32:
			var val float32
			err = binary.Read(tl.reader, endian, &val)
			tl.floatVal = float64(val)

		case TokenFloat64:
			err = binary.Read(tl.reader, endian, &tl.floatVal)
		}
	}

	if err != nil {
		return TokenEndOfSource, loc
	}

	return token, loc
}

// GetValueString is used to get a string value associated with tokens
// TokenString and TokenIdentifier.
func (tl *TokenList) GetValueString() string {
	return tl.strVal
}

func (tl *TokenList) GetValueInt64() int64 {
	return tl.intVal
}

func (tl *TokenList) GetValueUint64() uint64 {
	return tl.uintVal
}

func (tl *TokenList) GetValueFloat64() float64 {
	return tl.floatVal
}
