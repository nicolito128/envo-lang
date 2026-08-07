# The Envo Programming Language

Envo is an esoteric and minimalistic message-oriented programming language where all computing is expressed through message passing.

## Getting started

### Prerequisites 

* Go 1.26+

### Installation and execution

First clone the repository and run the interactive REPL:

```bash
    git clone https://n128.xyz/n128/envo-lang.git
    cd envo-lang
    go run ./cmd/envo/main.go
```

To build a standalone binary:

```bash
    go build -o envo ./cmd/envo/main.go
    ./envo
```

## Specification

For a formal description of the language grammar and AST node definitions, see the [SPEC](./SPEC.md).

## Language Tour

### 1. Messages

In Envo, the fundamental unit of structure is the **message**.

An **empty message** or **raw message** is represented with curly braces:

```envo
{}
```

Messages can contain multiple elements separated by commas:

```envo
{ 1, 2, ~ok, "hello" }
```

A **labeled message** attaches an identifier to a message body:

```envo
foo{}
digits{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
user{"name", "john"}
```

### 2. Data Types & Literals

| Type | Example Syntax |
| :--- | :--- |
| **Integers** | { 123, 0b1010, 0o755, 0xFF, 123e+3 } |
| **Floating-Point** | { 3.1415926, 1.2e-7, 0xFFp-4 } |
| **Strings** | { "hello, world\n", \`raw string without escape interpretation\` } |
| **Characters** | { 'a', '\n', 'Z' } |
| **Booleans** | { true, false }` |
| **Symbols** | { ~cat, ~dog, ~ok, ~error, ~_maybe } |

### 3. Definitions & Pattern Matching

You can bind behavior to a labeled message by specifying a **receiver**, a **pattern**, and a **response**:

```envo
receiver { pattern } { response }
```

#### Identity Rule
Creates an identity pattern that returns whatever argument it matches:

```envo
id{x}{x}
```

#### Factorial
Multiple definitions for the same receiver act as pattern-matched rules evaluated in sequence:

```envo
fact{0}{1}
fact{n}{n * fact{n - 1}}
```

#### Sum Up To N
```envo
sum_from{1}{1}
sum_from{x}{x + sum_from{x - 1}}
```
