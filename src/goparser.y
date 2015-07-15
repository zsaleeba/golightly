%{
	#include <stdio.h>
	#include "ast.h"
	#include "gl.h"
	#include "goparser.tab.h"
	
    extern int yylex();
    void yyerror(const char *s) { fprintf(stderr, "%d-%d:%d-%d:%s\n", yylloc.first_line, yylloc.last_line, yylloc.first_column, yylloc.last_column, s); }
%}

%defines
%locations

/* Represents the many different ways we can access our data */
%union {
    int        token;
	char      *ident;
	GoAst     *ast;
	int        intLit;
	double     floatLit;
	char      *stringLit;
	GoAstList *astList;
	GoAstCall *call;
}

/* Define our terminal symbols (tokens). This should
   match our tokens.l lex file. We also define the node type
   they represent.
 */
%token <intLit>        TINTEGER
%token <floatLit>      TDOUBLE
%token <stringLit>     TSTRING
%token <ident>         TIDENTIFIER
%token <token>         TLEFTBRACKET TRIGHTBRACKET TCOMMA TDOT
%token <token>         TPLUS TMINUS TMUL TDIV

/* Define the type of node our nonterminal symbols represent.
   The types refer to the %union declaration above. Ex: when
   we call an ident (defined by union type ident) we are really
   calling an (NIdentifier*). It makes the compiler happy.
 */
%type <ast>     stmt numeric string expr
%type <astList> program top_level call_args
%type <call>    func_call
%type <ident>   ident

/* Operator precedence for mathematical operators */
%left TPLUS TMINUS
%left TMUL TDIV

%start program

%%

program : 
        | top_level               { $$ = $1; GoRunProgram((GoAst *)$1) }
        ;
		
top_level : stmt                  { $$ = GoAstListCreate(GoAstTypeBlock); GoAstListAppend($$, $1); GoRunStatement((GoAst *)$1) }
      | top_level stmt            { $$ = GoAstListAppend($1, $2); GoRunStatement((GoAst *)$2) }
      ;

/*        
stmts : stmt                      { $$ = GoAstListCreate(GoAstTypeBlock); GoAstListAppend($$, $1) }
      | stmts stmt                { $$ = GoAstListAppend($1, $2) }
      ;
*/

stmt : func_call                  { $$ = (GoAst *)$1 }
     ;

func_call : ident TLEFTBRACKET call_args TRIGHTBRACKET  
                                  { $$ = GoAstCallCreate($1, $3) }
          ;
    
call_args : /*blank*/             { $$ = GoAstListCreate(GoAstTypeParamList) }
          | expr                  { $$ = GoAstListCreate(GoAstTypeParamList); GoAstListAppend($$, $1) }
          | call_args TCOMMA expr { $$ = GoAstListAppend($1, $3) }
          ;

ident : TIDENTIFIER
      ;

numeric : TINTEGER                { $$ = (GoAst *)GoAstIntLiteralCreate($1); }
        | TDOUBLE                 { $$ = (GoAst *)GoAstFloatLiteralCreate($1); }
        ;
    
string  : TSTRING                 { $$ = (GoAst *)GoAstStringLiteralCreate($1); }
        ;
    
expr : numeric
     | string
     ;

%%
