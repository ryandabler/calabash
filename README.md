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

BLOCK_STATEMENT
    : '{' PROGRAM '}'
    ;

MULTI_IDENT
    : MULTI_IDENT ',' IDENTIFIER
    | IDENTIFIER
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
    | FUNDAMENTAL
    ;

FUNDAMENTAL
    : NUMBER
    | STRING
    | '(' EXPRESSION ')'
    | 'bottom'
    | 'true'
    | 'false'
    ;
```
