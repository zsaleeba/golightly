/*
** $Id: llex.h,v 1.72 2011/11/30 12:43:51 roberto Exp $
** Lexical Analyzer
** See Copyright Notice in golightly.h
*/

#ifndef llex_h
#define llex_h

#include "lobject.h"
#include "lzio.h"


#define FIRST_RESERVED	257



/*
* WARNING: if you change the order of this enumeration,
* grep "ORDER RESERVED"
*/
enum RESERVED {
  /* terminal symbols denoted by reserved words */
  TK_AND = FIRST_RESERVED, TK_BREAK, TK_CASE, TK_CHAN, TK_CONST,
  TK_CONTINUE, TK_DEFAULT, TK_DEFER, TK_DO, TK_ELSE, TK_ELSEIF, 
  TK_END, TK_FALSE, TK_FALLTHROUGH, TK_FOR, TK_FUNC, TK_FUNCTION,
  TK_GO, TK_GOTO, TK_IF, TK_IMPORT, TK_IN, TK_INTERFACE, TK_LOCAL, 
  TK_MAP, TK_NIL, TK_NOT, TK_OR, TK_PACKAGE, TK_RANGE, TK_REPEAT,
  TK_RETURN, TK_SELECT, TK_STRUCT, TK_SWITCH, TK_THEN, TK_TRUE, 
  TK_TYPE, TK_UNTIL, TK_VAR, TK_WHILE,
  /* other terminal symbols */
  TK_CONCAT, TK_DOTS, TK_EQ, TK_GE, TK_LE, TK_NE, TK_DBCOLON, 
  TK_SHIFT_LEFT, TK_SHIFT_RIGHT, TK_BIT_CLEAR,
  TK_ADD_ASSIGN, TK_SUBTRACT_ASSIGN, TK_MULTIPLY_ASSIGN,
  TK_DIVIDE_ASSIGN, TK_MOD_ASSIGN,
  TK_BITWISE_AND_ASSIGN, TK_BITWISE_OR_ASSIGN, TK_BITWISE_EXOR_ASSIGN,
  TK_SHIFT_LEFT_ASSIGN, TK_SHIFT_RIGHT_ASSIGN, TK_BIT_CLEAR_ASSIGN,
  TK_LOGICAL_AND, TK_LOGICAL_OR, TK_CHANNEL_COMM, 
  TK_INCREMENT, TK_DECREMENT,
  TK_NOT_EQUAL, TK_DECLARE_ASSIGN,
  TK_EOS,
  TK_NUMBER, TK_NAME, TK_STRING
};

/* number of reserved words */
#define NUM_RESERVED	(cast(int, TK_WHILE-FIRST_RESERVED+1))


typedef union {
  lua_Number r;
  TString *ts;
} SemInfo;  /* semantics information */


typedef struct Token {
  int token;
  SemInfo seminfo;
} Token;


/* state of the lexer plus state of the parser when shared by all
   functions */
typedef struct LexState {
  int current;  /* current character (charint) */
  int linenumber;  /* input line counter */
  int lastline;  /* line of last token `consumed' */
  Token t;  /* current token */
  Token lookahead;  /* look ahead token */
  struct FuncState *fs;  /* current function (parser) */
  struct lua_State *L;
  ZIO *z;  /* input stream */
  Mbuffer *buff;  /* buffer for tokens */
  struct Dyndata *dyd;  /* dynamic structures used by the parser */
  TString *source;  /* current source name */
  TString *envn;  /* environment variable name */
  char decpoint;  /* locale decimal point */
} LexState;


LUAI_FUNC void luaX_init (lua_State *L);
LUAI_FUNC void luaX_setinput (lua_State *L, LexState *ls, ZIO *z,
                              TString *source, int firstchar);
LUAI_FUNC TString *luaX_newstring (LexState *ls, const char *str, size_t l);
LUAI_FUNC void luaX_next (LexState *ls);
LUAI_FUNC int luaX_lookahead (LexState *ls);
LUAI_FUNC l_noret luaX_syntaxerror (LexState *ls, const char *s);
LUAI_FUNC const char *luaX_token2str (LexState *ls, int token);


#endif
