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
    | WHILE
    | 'continue' ';'
    | 'break' ';'
    ;

VARIABLE_DECLARATION
    : 'let' MULTI_IDENT_DECL '=' MULTI_EXPR ';'
    ;

ASSIGNMENT
    : MULTI_IDENT '=' MULTI_EXPR ';'
    ;

MULTI_IDENT
    : MULTI_IDENT ',' identifier
    | identifier
    ;

IF
    : 'if' VARIABLE_DECLARATION? EXPRESSION BLOCK_STATEMENT ['else' (IF | BLOCK_STATEMENT)]?
    ;

WHILE
    : 'while' VARIABLE_DECLARATION? EXPRESSION BLOCK_STATEMENT
    ;

RETURN
    : 'return' EXPRESSION? ';'

BLOCK_STATEMENT
    : '{' PROGRAM '}'
    ;

MULTI_IDENT_DECL
    : MULTI_IDENT_DECL ',' IDENT_DECL
    | IDENT_DECL
    ;

IDENT_DECL
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
    | CALL_OR_GET
    ;

CALL_OR_GET
    : FUNDAMENTAL ('(' ARGUMENTS_LIST? ')' | '->' FUNDAMENTAL)*
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
    | PROTOTYPE
    | 'me'
    | RECORD
    ;

FUNCTION
    : 'fn' '(' MULTI_IDENT_DECL? ')' FUNC_BODY
    ;

TUPLE
    : '[' ARGUMENTS_LIST? ']'
    ;

FUNC_BODY
    : '->' EXPRESSION
    | BLOCK_STATEMENT
    ;

PROTOTYPE
    : 'proto' '{' PROTO_METHODS '}'
    ;

PROTO_METHODS
    : PROTO_METHODS ',' PROTO_METHOD
    | PROTO_METHOD
    ;

PROTO_METHOD
    : FUNDAMENTAL '->' FUNCTION
    ;

RECORD
    : '{' RECORD_KEY_VALUES? '}'
    ;

RECORD_KEY_VALUES
    : RECORD_KEY_VALUES ',' RECORD_KEY_VALUE
    | RECORD_KEY_VALUE
    ;

RECORD_KEY_VALUE
    : FUNDAMENTAL '->' FUNDAMENTAL
    ;
```
