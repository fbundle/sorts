## EL — a tiny typed, expression language

EL is a small, homoiconic, s‑expression language with a minimal core and a typed evaluation model built on the `sorts` library. Code and data share the same shape: lists and terms. The evaluator performs stepwise resolution under a persistent frame, tracking both inferred sorts (types) and values.

### Features
- **Homoiconic forms**: programs are made of terms and lists.
- **Function application**: simple two‑element list `(f x)`.
- **Lambdas**: `=>` special form.
- **Let‑bindings**: `let` special form for introducing names with declared sorts.
- **Pattern matching**: `match` with exact and structural patterns.
- **Type constructors**: function arrow `->`, sum `⊕`, product `⊗`.
- **Universes and built‑ins**: `U_n`, `Any_n`, `Unit_n` provide universe levels and terminal/initial sorts.

### Project layout
- `el/` — language core (AST, parser, resolver, matcher)
- `cmd/el/` — CLI that reads EL source from stdin
- `form/`, `sorts/` — shared form and sort systems used by EL

## Install / Build

### Requirements
- Go 1.22+

### Build the CLI
```bash
cd /Users/khanh/vault/private/code/sorts
go build -o ./bin/el ./cmd/el
```

Or install into your `GOBIN`:
```bash
go install ./cmd/el
```

Verify:
```bash
echo "(let x U_1 undef x)" | ./bin/el
```

## Language overview

### Syntax
- **Terms**: bare identifiers like `Nat`, `succ`, `x`.
- **Lists**: `(head arg)`; the empty list is not allowed as a term.
- **Application**: a regular two‑element list. Multi‑arg functions are curried: `((add 2) 3)`.

### Special forms

- **Lambda**: `(=> param body)`
  - Produces a function value without immediately evaluating the body.

- **Let**: `(let name1 type1 value1 ... nameN typeN valueN final)`
  - Binds a sequence of names, each with a declared sort and a value; evaluates `final` in the extended frame.
  - Use the literal `undef` to declare a name without assigning a value yet; the name becomes a self‑reference (useful for recursive references and data declarations).

- **Match**: `(match cond comp1 value1 ... compN valueN final)`
  - Tries each pattern `compK` against `cond`. On first match, evaluates `valueK` in the possibly extended frame (pattern variables are bound). If none match, evaluates `final`.

### Patterns
- **Exact**: `(exact expr)` matches only when the resolved value is syntactically equal to `expr`.
- **Variable**: a bare term in pattern position binds to the matched value (and its sort).
- **Structural**: `(f x)` patterns match function calls and recurse over command and argument.

### Sorts (types)
- **Arrow**: `(-> A B)` is the sort of functions from `A` to `B`.
- **Sum**: `(⊕ A B)` disjoint union of sorts.
- **Product**: `(⊗ A B)` pair/product of sorts.
- **Universes**: `U_n` denotes universe level `n+1` of sorts; `Any_n` is the terminal sort at level `n+1`; `Unit_n` is the initial sort at level `n+1`.

Type checking currently validates function application shapes against arrows; additional checks inside bindings and matches are being expanded.

## Example

`el/example.el` demonstrates natural numbers, addition, and simple matching:

```el
(let
    Bool U_1                undef
    True Bool               undef
    False Bool              undef

    Nat U_1                 undef
    0 Nat                   undef
    succ {Nat -> Nat}       undef

    1 Nat (succ 0)
    2 Nat (succ 1)
    3 Nat (succ 2)
    4 Nat (succ 3)

    x Any_0 {1 ⊕ 2 ⊕ 3}
    x Any_0 {1 ⊗ 2 ⊗ 3 ⊗ 4}

    is_two {Nat -> Bool} {x => (match x
        (exact (succ 1)) True
        False
    )}

    add {Nat -> Nat -> Nat} {x => {y => (match y
        (succ z) (succ ((add x) z))
        x
    )}}

    # (is_two 2)
    ((add 2) 3)             # output (succ (succ (succ (succ (succ 0)))))
)
```

Run it with the CLI:

```bash
./bin/el < /Users/khanh/vault/private/code/sorts/el/example.el
```

You will see, for each top‑level form, the printed form, its inferred sort, and its resolved value, e.g.:

```text
expr	 ((add 2) 3)
sort	 Nat
value	 (succ (succ (succ (succ (succ 0)))))
```

## CLI usage

- **Input**: reads UTF‑8 EL code from stdin; multiple top‑level forms are allowed one after another.
- **Output**: for each form, prints three lines: `expr`, `sort`, `value`.

Examples:
```bash
echo "(=> x x)" | ./bin/el
echo "(let id {Any_0 -> Any_0} {x => x} (id id))" | ./bin/el
```

## Notes

- Unicode operators `⊕` and `⊗` are regular identifiers; ensure your editor saves files as UTF‑8.
- EL uses explicit braces `{ ... }` in examples purely for readability; they are not required by the reader—any list is written with parentheses.
- The implementation favors clarity; sorts are reported using names from the `sorts` package.


