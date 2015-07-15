#include <stdio.h>
#include <jit/jit.h>
#include <jit/jit-dump.h>

#include "ast.h"
#include "codegen_jit.h"


void GoCodegenJitInit(GoCodegenJit *cj)
{
	cj->context = jit_context_create();
}

void GoCodegenJitClose(GoCodegenJit *cj)
{
	jit_context_destroy(cj->context);
}

int GoCodegenJitExecute(GoCodegenJit *cj, GoAst *ast)
{
	jit_context_build_start(cj->context);
	
	jit_function_t jitFunc = GoCodegenJitCompile(cj, ast);
	jit_dump_function(stdout, jitFunc, "func [uncompiled]");
	
    jit_function_compile(jitFunc);
	jit_context_build_end(cj->context);
 
    jit_dump_function(stdout, jitFunc, "func [compiled]");

	// Run it.
	jit_function_apply(jitFunc, NULL, NULL);
	
	return 0;
}


jit_function_t GoCodegenJitCompile(GoCodegenJit *cj, GoAst *ast)
{
	return NULL;
}
