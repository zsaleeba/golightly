%option outfile="golexer.yy.c"
%option yylineno

%{
#include "ast.h"
#include "goparser.tab.h"

#define TOKEN(t) (yylval.token = t)
#define SAVE_TOKEN 
int yywrap() { return 1; }
YYLTYPE yylloc;
int colnum = 1;

#define YY_USER_ACTION {                               \
		yylloc.first_line = yylineno;                  \
        yylloc.first_column = colnum;                  \
        colnum = colnum + yyleng;                      \
        yylloc.last_column = colnum;                   \
        yylloc.last_line = yylineno;                   \
	}
%}

%%

[ \t\n]                 ;
[a-zA-Z_][a-zA-Z0-9_]*  yylval.ident = strdup(yytext); return TIDENTIFIER;
[0-9]+\.[0-9]*          yylval.floatLit = strtod(yytext, NULL); return TDOUBLE;
[0-9]+                  yylval.intLit = atoi(yytext); return TINTEGER;
"\"[^\"]*\""            yylval.stringLit = strdup(yytext); return TSTRING;
"("                     return TOKEN(TLEFTBRACKET);
")"                     return TOKEN(TRIGHTBRACKET);
"."                     return TOKEN(TDOT);
","                     return TOKEN(TCOMMA);
.                       printf("Unknown token!\n"); yyterminate();

%%