# calabash

A simple hobby scripting language

## Grammar

The grammar for Calabash is:

```ebnf
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
    ;
```
