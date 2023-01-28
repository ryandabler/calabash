# calabash

A simple hobby scripting language

## Grammar

The grammar for Calabash is:

```ebnf
PROGRAM
    : (DECLARATION | EXPRESSION)*
    ;

DECLARATION
    : ASSIGNMENT
    ;

ASSIGNMENT
    : 'let' MULTI_IDENT '=' MULTI_EXPR ';'
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
    : ADDITION
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
