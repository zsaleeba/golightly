#ifndef GOLIGHTLY_AST_H
#define GOLIGHTLY_AST_H

typedef enum 
{
	GoAstTypeIntLiteral,
	GoAstTypeFloatLiteral,
	GoAstTypeStringLiteral,
	GoAstTypeCall,
	GoAstTypeBlock,
	GoAstTypeParamList,
	GoAstTypeList
} GoAstType;

typedef struct
{
	GoAstType nodeType;
} GoAst;

typedef struct
{
	GoAst      base;
	int        val;
} GoAstIntLiteral;

typedef struct
{
	GoAst      base;
	double     val;
} GoAstFloatLiteral;

typedef struct
{
	GoAst      base;
	const char *val;
} GoAstStringLiteral;

typedef struct
{
	int        used;
	int        capacity;
	GoAst    **asts;
} GoAstVec;

typedef struct
{
	GoAst      base;
	GoAstVec   vec;
} GoAstList;

typedef struct
{
	GoAst       base;
	const char *ident;
	GoAstVec    params;
} GoAstCall;

/* prototypes */
GoAstVec           *GoAstVecInit(GoAstVec *vec);
void                GoAstVecFree(GoAstVec *vec);
GoAstVec           *GoAstVecAppend(GoAstVec *vec, GoAst *node);
GoAstList          *GoAstListCreate(GoAstType typ);
void                GoAstListFree(GoAstList *node);
GoAstList          *GoAstListAppend(GoAstList *list, GoAst *node);
GoAstCall          *GoAstCallCreate(const char *ident, GoAstList *params);
void                GoAstCallFree(GoAstCall *node);
GoAstIntLiteral    *GoAstIntLiteralCreate(int val);
GoAstFloatLiteral  *GoAstFloatLiteralCreate(double val);
GoAstStringLiteral *GoAstStringLiteralCreate(const char *val);
GoAst              *GoAstSetType(GoAst *node, GoAstType typ);

#endif 
