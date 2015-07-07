#include <stdlib.h>
#include <string.h>

#include "ast.h"


#define NODE_LIST_INITIAL_SIZE 4


GoAstVec *GoAstVecInit(GoAstVec *vec)
{
	vec->used = 0;
	vec->capacity = NODE_LIST_INITIAL_SIZE;

	vec->asts = malloc(NODE_LIST_INITIAL_SIZE * sizeof(GoAst *));
	if (vec->asts == NULL)
	{
		free(vec->asts);
	    return NULL;
	}
	
	return vec;
}


void GoAstVecFree(GoAstVec *vec)
{
	if (vec->capacity > 0)
	{
		vec->used = 0;
		vec->capacity = 0;
		free(vec->asts);
	}
}


GoAstVec *GoAstVecAppend(GoAstVec *vec, GoAst *node)
{
	if (vec->used == vec->capacity)
	{
		/* Double the vector size. */
		vec->capacity *= 2;
		vec->asts = realloc(vec->asts, vec->capacity * sizeof(GoAst *));
		if (vec->asts == NULL)
			return NULL;
	}
	
	/* Append. */
	vec->asts[vec->used] = node;
	vec->used++;
	
	return vec;
}


GoAstList *GoAstListCreate(GoAstType typ)
{
	GoAstList *node = malloc(sizeof(GoAstList));
	if (node == NULL)
	    return NULL;

	node->base.nodeType = typ;
	if (GoAstVecInit(&node->vec) == NULL)
	{
		free(node);
		return NULL;
	}
	
	return node;
}


void GoAstListFree(GoAstList *node)
{
	GoAstVecFree(&node->vec);
	free(node);
}


GoAstList *GoAstListAppend(GoAstList *list, GoAst *node)
{
	if (GoAstVecAppend(&list->vec, node) == NULL)
		return NULL;
	
	return list;
}


GoAstCall *GoAstCallCreate(const char *ident, GoAstList *params)
{
	GoAstCall *node = malloc(sizeof(GoAstCall));
	if (node == NULL)
	    return NULL;

	node->base.nodeType = GoAstTypeCall;
	node->ident = ident;
	memcpy(&node->params, &params->vec, sizeof(params->vec));
	free(params);
	
	return node;
}


void GoAstCallFree(GoAstCall *node)
{
	GoAstVecFree(&node->params);
	if (node->ident != NULL)
		free((void *)node->ident);
		
	free(node);
}


GoAstIntLiteral *GoAstIntLiteralCreate(int val)
{
	GoAstIntLiteral *node = malloc(sizeof(GoAstIntLiteral));
	if (node == NULL)
	    return NULL;

	node->base.nodeType = GoAstTypeIntLiteral;
	node->val = val;
	
	return node;
}


GoAstFloatLiteral *GoAstFloatLiteralCreate(double val)
{
	GoAstFloatLiteral *node = malloc(sizeof(GoAstFloatLiteral));
	if (node == NULL)
	    return NULL;

	node->base.nodeType = GoAstTypeFloatLiteral;
	node->val = val;
	
	return node;
}


GoAstStringLiteral *GoAstStringLiteralCreate(const char *val)
{
	GoAstStringLiteral *node = malloc(sizeof(GoAstStringLiteral));
	if (node == NULL)
	    return NULL;

	node->base.nodeType = GoAstTypeStringLiteral;
	node->val = val;
	
	return node;
}


GoAst *GoAstSetType(GoAst *node, GoAstType typ)
{
	node->nodeType = typ;
	return node;
}
