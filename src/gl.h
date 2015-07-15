#ifndef GOLIGHTLY_GL_H
#define GOLIGHTLY_GL_H

#include "ast.h"


#ifndef TRUE
#define TRUE 1
#define FALSE 0
#endif

/* Prototypes. */
int yyparse();

int GoCodegenLlvm(GoAst *ast);

void GoRunProgram(GoAst *ast);
void GoRunStatement(GoAst *ast);

#endif /* GOLIGHTLY_GL_H */
