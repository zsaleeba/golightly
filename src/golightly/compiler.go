package golightly

type Compiler struct {
	files map[string]CompileFile
}

type CompileFile struct {
	filename string

}
