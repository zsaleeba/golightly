#ifndef GOLIGHTLY_CODEGEN_JIT_H
#define GOLIGHTLY_CODEGEN_JIT_H

#include <jit/jit.h>



typedef struct
{
	jit_context_t context;
} GoCodegenJit;


// Prototypes.
void GoCodegenJitInit(GoCodegenJit *cj);
void GoCodegenJitClose(GoCodegenJit *cj);
int GoCodegenJitExecute(GoCodegenJit *cj, GoAst *ast);
jit_function_t GoCodegenJitCompile(GoCodegenJit *cj, GoAst *ast);


#endif /* GOLIGHTLY_CODEGEN_JIT_H */
