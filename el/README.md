## EL — Language specification

EL is a small, homoiconic, typed s‑expression language. Code and data share the same structure; evaluation and typing are defined over a persistent frame using the `sorts` library.

### Syntax

- **Terms**: bare identifiers like `Nat`, `succ`, `x`.
- **Lists**: `(head arg)`; empty list is not a term.
- **Application**: two‑element list `(f x)`. Multi‑arg functions are curried: `((add 2) 3)`.

### Special forms

- **Lambda**: `(=> param body)`
  - Produces a function value; body is evaluated when applied.

- **Let**: `(let name1 type1 value1 ... nameN typeN valueN final)`
  - Binds each `nameK` with declared sort `typeK` and value `valueK` in order, then evaluates `final` in the extended frame.
  - Use literal `undef` to declare without assigning a value; the name becomes a self‑reference (enables recursive references and data declarations).

- **Match**: `(match cond comp1 value1 ... compN valueN final)`
  - Tries patterns `compK` against `cond` in order. On first match, evaluates `valueK` in the possibly extended frame (pattern variables are bound). If none match, evaluates `final`.

### Patterns

- **Exact**: `(exact expr)` matches only when the resolved value is syntactically equal to `expr`.
- **Variable**: a bare term in pattern position binds that name to the matched value and its sort.
- **Structural**: `(f x)` matches a function call and recurses on command and argument.

### Sorts (types)

- **Arrow**: `(-> A B)` functions from `A` to `B`.
- **Sum**: `(⊕ A B)` disjoint union of sorts.
- **Product**: `(⊗ A B)` pair/product of sorts.
- **Universes**: `U_n` is universe level `n+1`; `Any_n` terminal sort at level `n+1`; `Unit_n` initial sort at level `n+1`.

### Static typing

- **Function call**: `(f x)` checks iff `f : (-> A B)` and `x : A`, yielding result sort `B`.
- **Let‑binding**: `(let name T value ...)` checks each binding in order:
  - `undef` is always permitted.
  - lambdas require `T = (-> A B)`; extend with `name : T` and `param : A`, then check `body : B`.
  - `match` branches are checked by reverse pattern matching the case pattern against the scrutinee’s sort, then checking the branch against the expected sort; `final` is checked similarly.
  - other expressions resolve to a sort `S` that must satisfy `TermOf(S, T)`.

See `ALGORITHM.md` for precise matching, reverse matching, and typing procedures.

### Example

`example.el` demonstrates naturals, addition, and matching:

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

    is_two {Nat -> Bool} {x => (match x
        (exact (succ 1)) True
        False
    )}

    add {Nat -> Nat -> Nat} {x => {y => (match y
        (succ z) (succ ((add x) z))
        x
    )}}

    ((add 2) 3)
)
```

### Build and run

From the repository root:
```bash
go build -o ./bin/el ./cmd/el
./bin/el < el/example.el
```

### Notes

- Unicode operators `⊕` and `⊗` are regular identifiers (UTF‑8 source).
- Special‑form arities (`=>`, `let`, `match`, `->`, `⊕`, `⊗`, `exact`) are validated at parse time.
- Error messages point to ill‑typed applications and malformed forms during resolution.


