### Few Design Notes of funlang

Funlang syntax is inspired heavily from golang.

#### Lexer
This project avoids using most of the compiler tools for generating lexer,
parser or code generators because I made this project for learning how
compilers work. So lexer is hand written and is capable of parsing utf-8
encoding (Thanks to golang). There is no lookahead.

#### Parser
Parser is based on recursive descent parsing strategy with single token
lookahead. Parser does nothing regarding semantic checks. It just generates
correct AST for a program. There is no type checking and symbol resolution
till now.

#### Semantic analysis
Semantic Analysis exists inside sema/ directory. I have implemented symbol
resolution with scopes.


