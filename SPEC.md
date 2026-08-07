# Envo Programming Language Specification

## Introduction

Envo is an esoteric and minimalist message-oriented programming language where all computing is expressed through message passing.

## Notation

The syntax is specified using a [variant](https://en.wikipedia.org/wiki/Wirth_syntax_notation) of Extended Backus-Naur Form (EBNF):

## Programs

### Basic principles of a program

An Envo program is a collection of statements, where each statement is either a definition or an expression.

    program             = { statement } .
    statement           = define | expr .

## Lexical elements

### Comments

Inline comments start with `#` and stop at the end of the line.

    # Inline comment

Multi-line comments start with `#{` and stop at the next `}#`.

    #{
 Multi-line
        comment
    }#

### Separators

You can use the symbol `;` as a message separator, but it is optional. A message can be separated by a newline or a semicolon.

### Operators

A list of valid operators in Envo is provided below. Operators are used to perform operations on operands.

    unary_op            = "+" | "-" | "!" .

    binary_op           = "+" | "-" | "*" | "/" | "%" | "==" | "!=" | "<" | "<=" | ">" | ">=" | "&&" | "||" .

### Valid characters

    newline             := { unicode code point U+000A }
    unicode_char        := { any arbitrary character, except newline }
    unicode_digit       := { any unicode code categorized as digit }
    unicode_letter      := { any unicode code categorized as letter }

    letter              = unicode_letter | "_" .
    decimal_digit       = "0" ... "9" .
    binary_digit        = "0" | "1" .
    octal_digit         = "0" … "7" .
    hex_digit           = "0" … "9" | "A" … "F" | "a" … "f" .

### Identifiers

An `identifier` name program expressions such as declarations, functions, variables and types.

    identifier = letter { letter | unicode_digit } .

### Symbol literals

A symbol is a literal constant with a name, used to represent fixed values.

    symbol_lit = "~" identifier .

### Character literals

A character literal is a single Unicode character enclosed in single quotes. It can be a Unicode character or an escaped character.

    char_lit            = "'" ( unicode_value ) "'" .

    unicode_value       = unicode_char | escaped_char .
    escaped_char        = `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) .

### String literals

A string literal is a sequence of Unicode characters enclosed in double quotes or backticks. It can be a raw string literal or an interpreted string literal.

    string_lit                  = raw_string_lit | interpreted_string_lit .

    raw_string_lit              = "`" { unicode_char | newline } "`" .
    interpreted_string_lit      = `"` { unicode_value } `"` .

### Integer literals

A integer literal is a sequence of digits that represents an integer value. It can be in decimal, binary, octal or hexadecimal format.

    int_lit             = decimal_lit | binary_lit | octal_lit | hex_lit .

    decimal_lit         = "0" | decimal_digits .
    binary_lit          = "0" ( "b" | "B" ) [ "_" ] binary_digits .
    octal_lit           = "0" ( "o" | "O" ) [ "_" ] octal_digits .
    hex_lit             = "0" ( "x" | "X" ) [ "_" ] hex_digits .

    decimal_digits      = decimal_digit { [ "_" ] decimal_digit } .
    binary_digits       = binary_digit { [ "_" ] binary_digit } .
    octal_digits        = octal_digit { [ "_" ] octal_digit } .
    hex_digits          = hex_digit { [ "_" ] hex_digit } .

### Floating-point literals

A floating-point literal is a sequence of digits that represents a floating-point value. It can be in decimal or hexadecimal format.

    float_lit               = decimal_float_lit | hex_float_lit

    decimal_float_lit       = decimal_digits "." decimal_digits [ decimal_exponent ] .
    decimal_exponent        = ( "e" | "E" ) [ "+" | "-" ] decimal_digits .

    hex_float_lit           = "0" ( "x" | "X" ) hex_digits [ hex_exponent ] .
    hex_exponent            = ( "p" | "P" ) [ "+" | "-" ] decimal_digits .

### Imaginary literals

An imaginary literal is a sequence of digits that represents an imaginary value. It can be in decimal or hexadecimal format, and it is denoted by the suffix "i".

    imaginary_lit = ( int_lit | float_lit ) "i" .

## Message

The core concept of Envo is the message, which is a collection of words enclosed in curly braces. A message can be empty or contain a list of words separated by commas.

Optionally, a message can have a label, which is an identifier that precedes the message. A label can be used to identify the message and its contents.

    message             = [ label ] raw_message .

    raw_message         = "{" [ word_list ] "}" .
    word_list           = word { "," word } .
    word                = expr .

    label               = identifier .

For example, the following are valid messages:

    {}
    foo{}
    bar{1, 2, 3}
    baz{ "hello", "world" }
    qux{ ~symbol, 42, [1, 2, 3], (~a, ~b) }

## Definition

A definition is a statement that defines a receiver, a pattern and a response. A receiver is a label that identifies the definition. A pattern is a message that specifies the input of the definition. A response is a message that specifies the output of the definition.

    define              = receiver pattern response .

    receiver            = label .
    pattern             = raw_message .
    response            = raw_message .

For example, the following is a valid definition:

    add{ x, y }{ x + y }
    double{ x }{ 2 * x }
    factorial{ 0 }{ 1 }
    factorial{ n }{ n * factorial{ n - 1 } }

Definitions allow you to create dispatch behaviors based on the pattern of the message received. The pattern can be matched against the message received, and if it matches, the response is executed.

## Expression

An expression specifies the computation of a value by applying operators and functions to operands. 

    expr                = binary_expr | unary_expr | operand .

    binary_expr         = operand binary_operator operand .
    unary_expr          = unary_operator operand .

For example, the following are valid expressions:

    1 + 2
    3 * (4 - 5)
    !true
    double{ add{1, 2} }

### Operands

    operand         = literal | identifier | message .

    literal         = basic_lit | identifier .
    basic_lit       = int_lit | float_lit | char_lit | string_lit | symbol_lit .
