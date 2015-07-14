#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <libgen.h>
#include <jit/jit.h>
#include <jit/jit-dump.h>

#include "gl.h"

/* Global. */
int compile = FALSE;
	

#if 0
// Returns an uncompiled jit_function_t
// Note that jit_adder is a jit_function_t that's passed into this builder.
jit_function_t build_hello_world(jit_context_t context) 
{
    // void foo()
    jit_function_t F = jit_function_create(context, jit_type_create_signature(jit_abi_cdecl, jit_type_void, NULL, 0, 1));

	// Prepare calling printf: create its signature
	jit_type_t type_cstring = jit_type_create_pointer(jit_type_sys_char, 1);
	jit_type_t printf_signature = jit_type_create_signature(jit_abi_vararg, jit_type_int, &type_cstring, 1, 1);
	
	// Create a global string.
	jit_value_t hello_str = jit_value_create_long_constant(F, type_cstring, (long)"hello world\n");
	
	// printf("hello")
	/* jit_value_t res = */ jit_insn_call_native(F, "printf", printf, printf_signature, &hello_str, 1, JIT_CALL_NOTHROW);

	return F;
}
#endif

void GoRunProgram(GoAst *ast)
{
	if (compile)
		GoAstPrint(ast, 0);
}


int main(int argc, char **argv)
{
	if (argc >= 1 && strcmp(basename(argv[0]), "glc") == 0)
		compile = TRUE;
		
	printf("%s - golightly\n", basename(argv[0]));
	
	yyparse();

#if 0
	jit_context_t context = jit_context_create();
	
	jit_context_build_start(context);
	jit_function_t hello_world_func = build_hello_world(context);
	
	jit_dump_function(stdout, hello_world_func, "hello_world_func [uncompiled]");
	
    jit_function_compile(hello_world_func);
	jit_context_build_end(context);
 
    jit_dump_function(stdout, hello_world_func, "hello_world_func [compiled]");

	jit_function_apply(hello_world_func, NULL, NULL);

	
	jit_context_destroy(context);
#endif
	
    return 0;
}
