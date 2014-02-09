package golightly

type SrcLoc struct {
	Line   int
	Column int
}

type SrcSpan struct {
	start SrcLoc
	end   SrcLoc
}
