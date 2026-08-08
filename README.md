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

To build a standalone binary and start the REPL:

```bash
go build -o envo ./cmd/envo/main.go
./envo
```

Or you can provide a valid envo file to execute:

```bash
./envo ./path/to/file
```

## Specification

For a formal description of the language grammar, see the [SPEC](./SPEC.md).

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
| **Integers** | { 123, 0b1010, 0o755, 0xFF } |
| **Floating-Point** | { 3.1415926, 1.2e-7, 0xFFp-4, 123e+4 } |
| **Strings** | { "hello, world\n", \`raw string without escape interpretation\` } |
| **Characters** | { 'a', '\n', 'Z' } |
| **Booleans** | { true, false }` |
| **Symbols** | { ~cat, ~dog, ~ok, ~error, ~_maybe } |

### 3. Definitions & Pattern Matching

You can bind behavior to a labeled message by specifying a **receiver**, a **pattern**, and a **response**:

```envo
receiver { pattern } { response }
```

For example, this definition creates an [identity](https://en.wikipedia.org/wiki/Identity_function) pattern that returns whatever argument it matches:

```envo
id{x}{x}
```

For more examples, see the [examples directory](./examples/).
