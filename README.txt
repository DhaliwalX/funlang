**EXPERIMENTS GOING ON**

(TODO: Improve README)

Funlang (name probably will change)
===============================================================================
Funlang is a programming language inspired from Go's syntax. Even its syntax is
more simpler than Go's syntax. The main focus of this project is not to create
another language but to learn how modern compilers use SSA form internally to
optimize code. The funlang compiler as of now does not target any hardware or
any interpreter. It will just compile the AST to SSA form and will **try** to
simplify it as much as possible.

Building
===============================================================================
This project requires latest Go compiler (>1.13) with modules support.

    go build

Trying the compiler
===============================================================================
As of now, compiler is lot verbose and will print lot of noise.

    ./funlang -input example.fun


References
===============================================================================
    1) http://www.dcs.gla.ac.uk/~jsinger/ssa.html
    2) https://dl.acm.org/citation.cfm?id=357071
    3) http://doi.acm.org/10.1145/115372.115320
    4) http://ssabook.gforge.inria.fr/latest/book.pdf

And, Go and LLVM source code.

LICENSE
===============================================================================
MIT

Copyright 2019 Pushpinder Singh

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
