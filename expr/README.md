## expr

A tiny expression AST package for a LISP-like syntax with optional infix blocks. It provides:

- Tokenizer: splits a source string into tokens with simple preprocessing
- Parser: builds an AST `Expr` from tokens, supporting both prefix and infix forms
- Marshaler: serializes an `Expr` back to a flat token stream

### Concepts

- **Token**: `string` alias representing lexical units
- **Expr**: AST node interface with `Marshal() []Token`
  - **Term**: leaf node for identifiers, literals, and punctuation
  - **Node**: list node, serialized as `( <arg1> <arg2> ... )`

### Special tokens

- Blocks: `(` `)` and infix blocks `{` `}`
- Infix operators: `+` `×` `=>` `:` `,`
  - `+`, `×`: left-associative inside `{ ... }`
  - `=>`, `:`, `,`: right-associative inside `{ ... }`
- Strings: delimited by `"`, escape with `\`
- Line comments: everything after `#` on a line is removed

### Grammar (informal)

- Prefix block: `( a b c )` → `Node{a,b,c}`
- Infix block: `{ a OP b OP c }` where all `OP` are the same
  - For `+`/`×`: `{1 + 2 + 3}` → `(+ (+ 1 2) 3)`
  - For `=>`/`:`: `{x => y => body}` → `(=> x (=> y body))`
  - For `,` (list): `{a, b, c}` → `(, a (, b c))`

### Quickstart

```go
package main

import (
    "fmt"
    "github.com/khanh/sorts/expr"
)

func main() {
    src := "(add 1 {2 + 3})"
    toks := expr.Tokenize(src)
    e, rest, err := expr.Parse(toks)
    if err != nil { panic(err) }
    if len(rest) != 0 { panic("unexpected trailing tokens") }

    // Marshal back to tokens
    out := e.Marshal()
    fmt.Println(out)
}
```

### Examples

- Tokenizing
```go
expr.Tokenize("{a, b, c} # list")
// → ["{", "a", ",", "b", ",", "c", "}"]
```

- Parsing prefix
```go
ex, _, _ := expr.Parse(expr.Tokenize("(f x y)"))
// ex = Node{Term("f"), Term("x"), Term("y")}
```

- Parsing infix (left associative)
```go
ex, _, _ := expr.Parse(expr.Tokenize("{1 + 2 + 3}"))
// ex marshals to: ( + ( + 1 2 ) 3 )
```

- Parsing infix (right associative: lambda)
```go
ex, _, _ := expr.Parse(expr.Tokenize("{x => y => (add x y)}"))
// ex marshals to: ( => x ( => y ( add x y ) ) )
```

- Parsing infix (right associative: list with comma)
```go
ex, _, _ := expr.Parse(expr.Tokenize("{a, b, c}"))
// ex marshals to: ( , a ( , b c ) )
```

- Parsing strings and nested functions
```go
ex, _, _ := expr.Parse(expr.Tokenize("(print \"hello, world\" (upper (concat \"a\" \"b\")))"))
// ex marshals to: ( print "hello, world" ( upper ( concat "a" "b" ) ) )
```

### Notes

- The parser enforces that all operators within a single infix block `{ ... }` are the same and that the number of arguments is odd (operand/operator alternation).
- `Marshal()` for `Node` always emits surrounding `(` `)`. `Term` marshals to a single token.
- Errors include empty token streams, malformed infix blocks, or mixed operators in an infix block.

### License

See repository license.