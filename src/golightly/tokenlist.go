package golightly

import (
	"bytes"
	"encoding/binary"
)

// an encoded, compact form of the tokens
type TokenList struct {
	last   SrcLoc
	tokens *bytes.Buffer
	reader *Reader
	strVal string
	intVal int64
	uintVal uint64

}

const TokenFlagInt16 = 0xfd
const TokenFlagInt32 = 0xfe
const TokenFlagInt64 = 0xff

func NewTokenList(filename string) *TokenList {
	tl := new(TokenList)
	tl.last.Line = 1
	tl.last.Column = 1
	tl.tokens = new(bytes.Buffer)

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
	tl.tokens.WriteString(str)
}

func (tl *TokenList) AddFloat(pos SrcLoc, val float64) {
	tl.EncodeLoc(pos)
	v32 := float32(val)
	if float64(v32) == val {
		// can be represented as a float32
		tl.tokens.WriteByte(byte(TokenFloat32))
		binary.Write(tl.tokens, binary.LittleEndian, v32)
	} else {
		// we need the full 64 bits
		tl.tokens.WriteByte(byte(TokenFloat64))
		binary.Write(tl.tokens, binary.LittleEndian, val)
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
		binary.Write(tl.tokens, binary.LittleEndian, byte(val))
	} else if val < 0x10000 {
		binary.Write(tl.tokens, binary.LittleEndian, byte(TokenFlagInt16))
		binary.Write(tl.tokens, binary.LittleEndian, uint16(val))
	} else if val < 0x100000000 {
		binary.Write(tl.tokens, binary.LittleEndian, byte(TokenFlagInt32))
		binary.Write(tl.tokens, binary.LittleEndian, uint32(val))
	} else {
		binary.Write(tl.tokens, binary.LittleEndian, byte(TokenFlagInt64))
		binary.Write(tl.tokens, binary.LittleEndian, val)
	}
}


func (tl *TokenList) StartReading() {
	tl.reader = bytes.NewReader(tl.tokens.Bytes())
}

func (tl *TokenList) GetToken(loc *SrcLoc) Token {
	tl.DecodeLoc(loc)
	b, err := tl.reader.ReadByte()
	if err != nil {
		return TokenEndOfSource
	}

	token := Token(b)
	if b < int(TokenString) {
		// keywords and operators have no value
		return token
	} else {
		switch token {
		// literals
		case TokenString:
			strLen, err := tl.DecodeUint64()
			buf := make([]byte, strLen)

			return token

		case TokenRune
		case TokenInt
		case TokenUint
		case TokenFloat32
		case TokenFloat64

		// identifiers
		case TokenIdentifier

	}
}
