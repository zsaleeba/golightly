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

// Equals compares two source spans.
func (ss SrcSpan) Equals(to SrcSpan) bool {
	return ss.start.Equals(to.start) && ss.end.Equals(to.end)
}

// Equals compares two source spans.
func (ss SrcLoc) Equals(to SrcLoc) bool {
	return ss.Line == to.Line && ss.Column == to.Column
}
