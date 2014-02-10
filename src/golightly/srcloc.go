package golightly

// type SrcLoc gives a location in the source file.
type SrcLoc struct {
	Line   int
	Column int
}

// type SrcSpan gives a from/to range in the source file.
type SrcSpan struct {
	start SrcLoc
	end   SrcLoc
}

// Add adds two source spans to make a wider span. They must be in order.
func (ss SrcSpan) Add(to SrcSpan) SrcSpan {
	return SrcSpan{ss.start, to.end}
}
