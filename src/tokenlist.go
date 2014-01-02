package golightly

// an encoded, compact form of the tokens
type TokenList struct {
	last   SrcLoc
	tokens []byte
}

const TokenListSizeStart = 256

const TokenFlagInt16 = 0xfd
const TokenFlagInt32 = 0xfe
const TokenFlagInt64 = 0xff

func NewTokenList(filename string) *TokenList {
	tl := new(TokenList)
	tl.last.Line = 1
	tl.last.Column = 1
	tl.tokens = make([]byte, 0, TokenListSizeStart)

	return tl
}

func (tl *TokenList) Add(pos SrcLoc, token int) {
	tl.EncodeLoc(pos)
	tl.EncodeUint64(uint64(token))
}

func (tl *TokenList) AddInt(pos SrcLoc, token int, val int64) {
	tl.Add(pos, token)
	tl.EncodeInt64(val)
}

func (tl *TokenList) AddUInt(pos SrcLoc, token int, val uint64) {
	tl.Add(pos, token)
	tl.EncodeUint64(val)
}

func (tl *TokenList) AddString(pos SrcLoc, token int, str string) {
	tl.Add(pos, token)
	tl.EncodeString(str)
}

// EncodeLoc stores the location of this token as a delta from the
// previous token. If this token is on the same line as the previous
// token the number of columns to advance is stored as a positive
// integer. If the line has changed the negative of the absolute
// column is stored and then the number of lines to advance is stored.
func (tl *TokenList) EncodeLoc(pos SrcLoc) {
	if tl.last.Line == pos.Line {
		tl.EncodeInt64(int64(pos.Column - tl.last.Column)) // positive
		tl.last.Column = pos.Column
	} else {
		tl.EncodeInt64(-int64(pos.Column)) // negative
		tl.EncodeUint64(uint64(pos.Line - tl.last.Line))
		tl.last = pos
	}
}

// EncodeString encodes a unicode string. It firstly stores the byte
// length of the string using EncodeUint64, then stores the contents of
// the string.
func (tl *TokenList) EncodeString(str string) {
	// encode the string length
	tl.EncodeUint64(uint64(len(str)))

	// add the string at the end
	tl.tokens = append(tl.tokens, []byte(str)...)
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
	// output a size flag if we need to
	if val >= TokenFlagInt16 {
		if val < 0x10000 {
			tl.EncodeByte(TokenFlagInt16)
		} else if val < 0x100000000 {
			tl.EncodeByte(TokenFlagInt32)
		} else {
			tl.EncodeByte(TokenFlagInt64)
		}
	}

	// output the value
	tl.EncodeByte(byte(val))
	if val >= TokenFlagInt16 {
		tl.EncodeByte(byte(val >> 8))
		if val > 0x10000 {
			tl.EncodeByte(byte(val >> 16))
			tl.EncodeByte(byte(val >> 24))
			if val > 0x100000000 {
				tl.EncodeByte(byte(val >> 32))
				tl.EncodeByte(byte(val >> 40))
				tl.EncodeByte(byte(val >> 48))
				tl.EncodeByte(byte(val >> 56))
			}
		}
	}
}

// EncodeByte adds a byte to the token buffer
func (tl *TokenList) EncodeByte(val byte) {
	if len(tl.tokens) >= cap(tl.tokens) {
		// increase the space available
		newTokens := make([]byte, len(tl.tokens), cap(tl.tokens)*2)
		copy(newTokens, tl.tokens)
		tl.tokens = newTokens
	}

	// add the byte
	tl.tokens = append(tl.tokens, val)
}

// Compact reduces the memory used by the token buffer to avoid wastage
func (tl *TokenList) Compact() {
	newTokens := make([]byte, len(tl.tokens))
	copy(newTokens, tl.tokens)
	tl.tokens = newTokens
}
