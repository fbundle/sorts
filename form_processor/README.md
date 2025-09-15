## form — minimal s‑expression and infix forms

`form` provides two things:

- **Tokenization** from a source string into a flat stream of tokens
- **Parsing** tokens into a tiny AST with two node kinds: `Term` and `List`

It supports both prefix lists, delimited by `(` `)`, and a lightweight infix block, delimited by `{` `}` with a fixed set of associative operators. Strings, whitespace, and line comments are handled by the tokenizer.

---

### Install

```bash
go get github.com/fbundle/el_sorts/form_processor
```

---

### Quick start

```go
package main

import (
    "fmt"
    "strings"
    "github.com/fbundle/sorts/form"
)

func main() {
    src := `
        (add 1 2)
        {x => y => (add x y)}
        {1 + 2 + 3}
    `

    tokens := form.Tokenize(src)

    // Parse all top-level forms in sequence
    for len(tokens) > 0 {
        node, rest, err := form.Parse(tokens)
        if err != nil { panic(err) }
        tokens = rest

        // Round-trip back to tokens
        fmt.Println(strings.Join(node.Marshal(), " "))
    }
}
```

Output:

```text
( add 1 2 )
( => x ( => y ( add x y ) ) )
( + ( + 1 2 ) 3 )
```

---

### Concepts and data model

- **Token**: `type Token = string`
- **Form**: the parsed node interface
  - `Term string` — an atom/identifier/number or delimiter captured as one token
  - `List []Form` — a parenthesized list `( … )` or the normalized result of an infix expression

API surface:

```go
// Tokenization
func Tokenize(s string) []Token

// Parsing one form_processor from the front of a token slice
func Parse(tokenList []Token) (Form, []Token, error)

// AST nodes implement:
type Form interface {
    Marshal() []Token
}
type Term string
type List []Form
```

`Parse` returns the parsed `Form`, the remaining `[]Token` (so you can parse multiple top-level forms), and an `error`.

`Marshal` serializes a node back into a flat token stream. For `List`, it includes surrounding `(` and `)`.

---

### Syntax supported by the tokenizer

- **Lists**: `(` begins, `)` ends
- **Infix blocks**: `{` begins, `}` ends (see next section for operators)
- **Strings**: delimited by `"` … `"`, with escape `\` inside strings
- **Whitespace**: any Unicode whitespace separates tokens
- **Line comments**: `#` to end-of-line is removed before tokenization

Examples:

```text
(greet "hello, world")     # comment after a form
{x -> y -> z}               # right-associative arrows
{1 + 2 + 3}                 # left-associative sum
"a \"quoted\" word"        # strings keep escapes
```

Note: the tokenizer recognizes the special token `$` for future features; it is not interpreted by the parser in this package.

---

### Infix blocks `{ … }`

Inside `{ … }` the parser normalizes a flat, odd-length sequence of alternating operands and a repeated operator into a prefix `List` with explicit associativity. All operators within a single infix block must be the same; mixing operators yields an error.

- **Left-to-right associative operators**:
  - `+` (sum)
  - `*` (product)
- **Right-to-left associative operators**:
  - `⊕` (sum)
  - `⊗` (product)
  - `=>` (lambda-like)
  - `->` (arrow type)
  - `:` (type ascription)
  - `,` (comma list)
  - `=` (equality)
  - `:=` (binding)

Normalization examples:

```text
{1 + 2 + 3}            =>  ( + ( + 1 2 ) 3 )
{x => y => f x y}      =>  ( => x ( => y ( f x y ) ) )
{a , b , c}            =>  ( , a ( , b c ) )
```

Errors raised by infix parsing:

- Empty block: parsed as an empty list `()`
- Single element: returned as-is
- Even-length sequence: error "infix syntax must have an odd number of arguments"
- Operator position not a `Term`: error "infix operator must be a term"
- Mixed operators: error "infix operator must be the same <op>"


---

### Parsing strategy and streaming

`Parse` consumes from the front of a token slice and returns the remainder, enabling a simple loop to parse multiple top-level forms from one input. For performance, `Tokenize` pre-sorts known split tokens so longer tokens that are prefixes of shorter ones (e.g. `:=` vs `:`) are matched first.

---

### Round-tripping and printing

`Marshal()` provides a token-level round-trip. To render a human-readable string, join with spaces:

```go
func render(f form.Form) string {
    return strings.Join(f.Marshal(), " ")
}
```

Example:

```go
f, rest, _ := form.Parse(form.Tokenize("{1 * 2 * 3}"))
_ = rest // empty
fmt.Println(strings.Join(f.Marshal(), " "))
// ( * ( * 1 2 ) 3 )
```

---

### Limits and non-goals

- Only a fixed set of infix operators is recognized. Extend by editing `defaultParser.Split` and the associativity in `processInfix` in `parser.go`.
- This package does not perform evaluation or type checking—only tokenization and shape normalization.
- No reader macros beyond `{ … }` infix blocks.

---

### Testing ideas

- Tokenization of strings with escapes and comments
- Infix normalization associativity and error cases
- Round-trip `Marshal()` invariants for `Term` and `List`

---

### See also

- The higher-level `el` package that builds on these forms lives in `el/`.


