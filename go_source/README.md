# SORTS
started from some basic dependent type in [lean4](https://github.com/leanprover/lean4), I stated this project with the goal to implement the full dependent type system so that it is capable for mathemtical proof. This is probably a decade-long project, hope it would last.


This repository currently focuses on:
- Basic sort universe machinery with levels and a partial order (subtyping/cumulativity)
- Core type constructors: arrows (function types) and products (pairs)
- A tiny s-expression based front-end (EL) and a Python helper to transpile a Python-like syntax into EL

The vision is to iteratively extend this core into a dependable, minimal kernel with an ergonomic front-end suitable for mathematical proofs.

## Quick start

- Ensure Go is installed (see `go.mod` for version).
- Run the EL interpreter on the provided examples:

```bash
cat example.el | go run cmd/el/main.go
python cmd/elt/elt.py example.elt | go run cmd/el/main.go
```

## Project layout

- `sorts/`: core sort/type constructors and subtyping
  - `atom.go`: atomic sorts and universe chaining
  - `arrow.go`: function types, intro/elim
  - `prod.go`: product types, intro/elim
  - `sorts.go`: shared interfaces and helpers
- `universe/`: universe implementation and builtin rules (initial/terminal chains, cumulativity)
- `form/`, `form_processor/`: tokenization and parsing of s-expressions used by EL
- `cmd/el/`: CLI for interpreting EL input
- `cmd/elt/`: Python helper to transpile a lightweight Python-like syntax into EL
- `example.el`, `example.py`: sample programs

## Current capabilities (work-in-progress)

- Universe levels with initial/terminal chains and basic cumulativity checks
- Subtyping via a partial order on sorts, including function-type variance
- Intro/elim rules for arrows and products (at the sort level)
- EL front-end to construct and inspect terms; Python-like transpiler for convenience

## Syntax (EL and ELT)

EL is a tiny s-expression based language with optional infix sugar. Parsing happens in two layers:

- Tokens: split on `+`, `*`, `$`, `⊕`, `⊗`, `Π`, `Σ`, `=>`, `->`, `:`, `,`, `=`, `:=`.
- Blocks:
  - `(` … `)` builds lists directly: `(head arg1 arg2)`.
  - `{` … `}` is infix sugar re-associated into prefix form by precedence/associativity:
    - Left-assoc: `{1 + 2 + 3}` → `(+ (+ 1 2) 3)`; similarly for `*`.
    - Right-assoc (default): `{x => y => (add x y)}` → `(=> x (=> y (add x y)))`, `{A -> B -> C}` → `(-> A (-> B C))`.

Core EL forms handled by the compiler:

- Arrows: `(-> A B)` is the function sort from `A` to `B`. Infix sugar: `{A -> B}`.
- Products: `(⊗ A B)` is the product sort; sum `(⊕ A B)` is available at the sort level.
- Lambda: `(=> (: x A) body)` creates a lambda from `A` to `type(body)`. With infix sugar: `{x => body}` nests right-associatively for multiple params.
- Application: `(cmd arg)` is application; by default, lists that do not match a declared head are compiled as application (β-site), with type checking against the arrow type of `cmd`.
- Let-binding: `(let (:= x v1) (:= y v2) final)` sequentially binds names, then compiles `final` in the extended context.
- Match: `(match cond (=> pat1 v1) ... (=> _ vDefault))` computes the least upper bound of branch result sorts for the overall type.
- Inhabitant: `(inh T)` creates a fresh placeholder term of sort `T`.
- Inspect: `(inspect v)` prints a debug view and returns `v` unchanged.

ELT is a thin Python-like surface that transpiles to EL via `cmd/elt/elt.py`. Use it for ergonomics; semantics are identical after transpilation.

Small examples:

```el
{A -> B -> A}                   # right-assoc arrow
(=> (: x A) (=> (: y B) x))     # lambda
(let (:= id (=> (: x A) x)) (inspect id))
```

## Semantics

- Sorts and universes:
  - Every term is a `Sort` with attributes: concrete `Form`, `level` (universe level), `parent` (its type/sort), and `lessEqual` (subtyping/cumulativity predicate).
  - The universe provides initial/terminal chains per level, with cumulativity (e.g., `A ≤ Any` in a given level chain).
- Subtyping/partial order:
  - Functions use standard variance: `(A1 -> B1) ≤ (A2 -> B2)` when `A2 ≤ A1` and `B1 ≤ B2`.
  - Products/sums compute parents and levels componentwise, and participate in LUB/GLB where defined.
- Contexts and names:
  - A `Context` maps `Name → Sort`, layered over builtins from the universe (e.g., `Unit`, `Any`).
  - Setting a binding extends the ordered map; lookup resolves user scope first, then builtins.
- Typability and parents:
  - Compiling a form yields a `Sort` whose `parent` is its type; printing in compile mode shows `(type T - level n)`. In eval/debug modes, the concrete form is also shown.
- Reducible forms:
  - β (application), `let`, and `match` are represented as reducible nodes. In `ModeEval`, they will step (evaluation hooks are scaffolded and will be extended); in `ModeComp`, they remain as typed nodes.

## Compiler/interpreter algorithm

End-to-end pipeline used by `cmd/el`:

1. Tokenize: raw text → tokens, with comments stripped (lines starting with `#`).
2. Parse: recursive descent with block handling to produce `form.Form` AST nodes: `Name` or `List`.
3. Compile: walk the AST under a `Context`:
   - If a list head matches a registered compiler (`->`, `⊕`, `⊗`, `=>`, `inh`, `let`, `match`, `inspect`), dispatch to the corresponding constructor/typechecker.
   - Otherwise treat the list as application and enforce that the head’s parent is an arrow and the argument is a subtype of the arrow domain.
   - Each constructor computes the resulting node’s `parent` (its type) and `level`, and may extend the context (e.g., lambda parameter, let-bindings) before compiling subterms.
4. Print: the CLI prints each compiled node using `ctx.ToString`.

Modes:

- Compile mode (`ModeComp`, default): type information only, prints `(type T - level n)`.
- Eval mode (`ModeEval`): reducible nodes call `Reduce` to step (β, let, match; currently WIP).
- Debug mode (`ModeDebug`): like eval but emits per-node compile logs.

Error handling:

- Name resolution errors panic with `name_not_found`.
- Ill-typed constructs panic with `type_check`/`type_err` depending on stage.
- Parser enforces infix well-formedness and block balancing.

## Roadmap toward a dependent type theory proof assistant

Short-to-medium term:
- Add sum types and unit/empty as proper inhabitants
- Dependent function types (Π-types) and dependent pairs (Σ-types)
- Judgmental equality/definitional equality and normalization for the kernel
- Pattern matching and inductive families (starting with simple inductives)
- Contexts, typing judgments, and well-scoped variables with hygienic names
- Improved parser and elaboration from surface syntax to core terms

Long-term milestones:
- Universe polymorphism and cumulativity scaling
- General inductive types (W-types), recursion/recursors, and termination checking
- Tactic framework and interactive proving experience
- Standard library of datatypes and theorems
- Proof irrelevance/prop hierarchy exploration and potential classical axioms as options

## Contributing

This is an early-stage exploration. Issues, ideas, and small PRs are welcome—especially around tests, examples, and small core extensions. If proposing larger changes (new connectives, typing rules, or parser features), please open an issue first to discuss the design.

## License

TBD. For now, assume all rights reserved by the author until a license is added.

