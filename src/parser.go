package golightly

type Parser struct {
	filename    string // the name of the current file being parsed
	tokenLoc    SrcLoc // the location in the file of the current token
	tokenEndLoc SrcLoc // the location in the fiule of the end of the current token

}
