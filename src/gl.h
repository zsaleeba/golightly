#ifndef GL_H
#define GL_H

#include "ast.h"

#ifndef TRUE
#define TRUE 1
#define FALSE 0
#endif

/* Globals. */
GoAst *program;

/* Prototypes. */
int yyparse();
int GoCodegenJit(GoAst *ast);
int GoCodegenLlvm(GoAst *ast);
void GoRunProgram(GoAst *ast);

#endif /* GL_H */
