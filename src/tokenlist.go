package golightly

// an encoded, compact form of the tokens
type TokenList struct {
	last SrcLoc
	tokens []byte
}

const TokenListSizeStart = 256

func NewTokenList(filename string) {
	tl := new(TokenList)
	tl.last.Line = 1
	tl.last.Column = 1
	tl.tokens = make(TokenList, 0, TokenListSizeStart)
	
	return tl
}

func (tl *TokenList) Add(pos SrcLoc, token int) {
	tl.EncodeLoc(pos)
	tl.EncodeInt(token)
}

func (tl *TokenList) AddInt(pos SrcLoc, token int, val int64) {
	tl.Add(pos, token)
	tl.EncodeInteger(val)
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
	if last.Line == pos.Line {
		tl.EncodeInt(pos.Column - last.Column)   // positive
		last.Column = pos.Column
	} else {
		tl.EncodeInt(-pos.Column)                // negative
		tl.EncodeInt(pos.Line - last.Line)
		last = pos
	}
}

func (tl *TokenList) EncodeInt(val int) {
	if val < 0 {
		if val == MIN_INT {
			// annoying edge case - MIN_INT can't be represented in positive form
			EncodeByte(0x81)
			for i:=0; i<8; i++ {
				EncodeByte(0x80)
			}
			EncodeByte(0x02)
		} else {
			// negative number
			EncodeUint64(-val << 1)
		}
	} else {
		// positive number
		EncodeUint64(val << 1)
	}
}

// EncodeUint64 encodes a number using a variable precision method.
// 
func (tl *TokenList) EncodeUint64(val uint64) {
	precision := 0
	if val > 0xfc {
		tl.EncodeByte(0xfc)
	} else if val > 0x3fff {
		tl.EncodeByte(0xfd)
		tl.EncodeByte(val & 0xff)
		tl.EncodeByte(val >> 8)
	} else if val > 0x3fffffff {
		tl.EncodeByte(0xfe)
	} else {
		tl.EncodeByte(0xff)
	}
}
