package golightly

struct Parser {
	filename string            // the name of the current file being parsed
	tokenLoc ParseLoc          // the location in the file of the current token
	tokenEndLoc ParseLoc       // the location in the fiule of the end of the current token
	
}

struct ParseLoc {
	Line int
	Column int
}
