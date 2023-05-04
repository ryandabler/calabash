# calabash

A simple hobby scripting language

## Grammar

The grammar for Calabash is:

```ebnf
PROGRAM
    : (STATEMENT | EXPRESSION)*
    ;

STATEMENT
    : VARIABLE_DECLARATION
    | ASSIGNMENT
    | IF
    | RETURN
    ;

VARIABLE_DECLARATION
    : 'let' MULTI_IDENT '=' MULTI_EXPR ';'
    ;

ASSIGNMENT
    : MULTI_IDENT '=' MULTI_EXPR ';'
    ;

IF
    : 'if' ASSIGNMENT? EXPRESSION BLOCK_STATEMENT ['else' (IF | BLOCK_STATEMENT)]?
    ;

RETURN
    : 'return' EXPRESSION? ';'

BLOCK_STATEMENT
    : '{' PROGRAM '}'
    ;

MULTI_IDENT
    : MULTI_IDENT ',' IDENT
    | IDENT
    ;

IDENT
    : 'mut'? identifier
    ;

MULTI_EXPR
    : MULTI_EXPR ',' EXPRESSION
    | EXPRESSION
    ;

EXPRESSION
    : BOOLEAN_OR
    ;

BOOLEAN_OR
    : BOOLEAN_OR '||' BOOLEAN_AND
    | BOOLEAN_AND
    ;

BOOLEAN_AND
    : BOOLEAN_AND '&&' EQUALITY
    | EQUALITY
    ;

EQUALITY
    : COMPARISON [('==' | '!=') COMPARISON]?
    ;

COMPARISON
    : ADDITION [('<' | '<=' | '>' | '>=') ADDITION]?
    ;

ADDITION
    : ADDITION ('+' | '-') MULTIPLICATION
    | MULTIPLICATION
    ;

MULTIPLICATION
    : MULTIPLICATION ('*' | '/') EXPONENTIATION
    | EXPONENTIATION
    ;

EXPONENTIATION
    : EXPONENTIATION '**' UNARY
    | UNARY
    ;

UNARY
    : '-' UNARY
    | CALL
    ;

CALL
    : CALL '(' ARGUMENTS_LIST? ')'
    | FUNDAMENTAL
    ;

ARGUMENTS_LIST
    : ARGUMENTS_LIST ',' EXPRESSION
    | EXPRESSION
    ;

FUNDAMENTAL
    : number
    | string
    | identifier
    | '(' EXPRESSION ')'
    | 'bottom'
    | 'true'
    | 'false'
    | FUNCTION
    | TUPLE
    ;

FUNCTION
    : 'fn' '(' MULTI_IDENT* ')' FUNC_BODY
    ;

TUPLE
    : '[' ARGUMENTS_LIST? ']'
    ;

FUNC_BODY
    : '->' EXPRESSION
    | BLOCK_STATEMENT
    ;
```
